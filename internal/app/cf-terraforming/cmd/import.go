package cmd

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	cfv0 "github.com/cloudflare/cloudflare-go"
	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"github.com/zclconf/go-cty/cty"
)

// resourceImportStringFormats contains a mapping of the resource type to the
// composite ID that is compatible with performing an import.
var resourceImportStringFormats = map[string]string{
	"cloudflare_access_application":    ":account_id/:id",
	"cloudflare_access_group":          ":account_id/:id",
	"cloudflare_access_rule":           ":identifier_type/:identifier_value/:id",
	"cloudflare_account_member":        ":account_id/:id",
	"cloudflare_argo":                  ":zone_id/argo",
	"cloudflare_bot_management":        ":zone_id",
	"cloudflare_byo_ip_prefix":         ":id",
	"cloudflare_certificate_pack":      ":zone_id/:id",
	"cloudflare_custom_hostname":       ":zone_id/:id",
	"cloudflare_custom_pages":          ":identifier_type/:identifier_value/:id",
	"cloudflare_custom_ssl":            ":zone_id/:id",
	"cloudflare_filter":                ":zone_id/:id",
	"cloudflare_firewall_rule":         ":zone_id/:id",
	"cloudflare_healthcheck":           ":zone_id/:id",
	"cloudflare_ip_list":               ":account_id/:id",
	"cloudflare_load_balancer":         ":zone_id/:id",
	"cloudflare_load_balancer_pool":    ":account_id/:id",
	"cloudflare_load_balancer_monitor": ":account_id/:id",
	"cloudflare_origin_ca_certificate": ":id",
	"cloudflare_page_rule":             ":zone_id/:id",
	"cloudflare_rate_limit":            ":zone_id/:id",
	"cloudflare_record":                ":zone_id/:id",
	"cloudflare_ruleset":               ":identifier_type/:identifier_value/:id",
	"cloudflare_spectrum_application":  ":zone_id/:id",
	"cloudflare_teams_list":            ":account_id/:id",
	"cloudflare_teams_location":        ":account_id/:id",
	"cloudflare_teams_proxy_endpoint":  ":account_id/:id",
	"cloudflare_teams_rule":            ":account_id/:id",
	"cloudflare_tunnel":                ":account_id/:id",
	"cloudflare_turnstile_widget":      ":account_id/:id",
	"cloudflare_waf_override":          ":zone_id/:id",
	"cloudflare_waiting_room":          ":zone_id/:id",
	"cloudflare_worker_route":          ":zone_id/:id",
	"cloudflare_workers_kv_namespace":  ":account_id/:id",
	"cloudflare_zone_lockdown":         ":zone_id/:id",
	"cloudflare_zone":                  ":id",
}

var providerVersionString string

func init() {
	rootCmd.AddCommand(importCommand)
}

var importCommand = &cobra.Command{
	Use:    "import",
	Short:  "Output `terraform import` compatible commands in order to import resources into state",
	Run:    runImport(),
	PreRun: sharedPreRun,
}

func runImport() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		zoneID = viper.GetString("zone")
		accountID = viper.GetString("account")
		workingDir := viper.GetString("terraform-install-path")
		execPath := viper.GetString("terraform-binary-path")

		// Download terraform if no existing binary was provided
		if execPath == "" {
			tmpDir, err := os.MkdirTemp("", "tfinstall")
			if err != nil {
				log.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			installConstraints, err := version.NewConstraint("~> 1.0")
			if err != nil {
				log.Fatal("failed to parse version constraints for installation version")
			}

			installer := &releases.LatestVersion{
				Product:     product.Terraform,
				Constraints: installConstraints,
			}

			execPath, err = installer.Install(context.Background())
			if err != nil {
				log.Fatalf("error installing Terraform: %s", err)
			}
		}

		// Setup and configure Terraform to operate in the temporary directory where
		// the provider is already configured.
		log.Debugf("initializing Terraform in %s", workingDir)
		tf, err := tfexec.NewTerraform(workingDir, execPath)
		if err != nil {
			log.Fatal(err)
		}

		_, providerVersion, err := tf.Version(context.Background(), true)
		providerVersionString = providerVersion[providerRegistryHostname+"/cloudflare/cloudflare"].String()
		log.WithFields(logrus.Fields{
			"version": providerVersionString,
		}).Debug("detected provider")

		var jsonStructData []interface{}

		if strings.HasPrefix(providerVersionString, "5") {
			resources := strings.Split(resourceType, ",")
			for _, resourceType := range resources {
				var result *http.Response

				// by default, we want to use the `list` operation however, there are times
				// when resources exist only as `get` operations but contain multiple
				// resources.
				endpoint := resourceToEndpoint[resourceType]["list"]
				if endpoint == "" {
					endpoint = resourceToEndpoint[resourceType]["get"]
				}

				// if we encounter a combined endpoint, we need to rewrite to use the correct
				// endpoint depending on what parameters are being provided.
				if strings.Contains(endpoint, "{account_or_zone}") {
					if accountID != "" {
						endpoint = strings.Replace(endpoint, "/{account_or_zone}/{account_or_zone_id}/", "/accounts/{account_id}/", 1)
					} else {
						endpoint = strings.Replace(endpoint, "/{account_or_zone}/{account_or_zone_id}/", "/zones/{zone_id}/", 1)
					}
				}

				// replace the URL placeholders with the actual values we have.
				placeholderReplacer := strings.NewReplacer("{account_id}", accountID, "{zone_id}", zoneID)
				endpoint = placeholderReplacer.Replace(endpoint)

				if strings.Contains(endpoint, "{") {
					log.WithFields(logrus.Fields{
						"resource": resourceType,
						"endpoint": endpoint,
					}).Debug("failed to substitute all path placeholders due to unknown parameters")

					continue
				}

				client := cloudflare.NewClient()

				err := client.Get(context.Background(), endpoint, nil, &result)
				if err != nil {
					var apierr *cloudflare.Error
					if errors.As(err, &apierr) {
						if apierr.StatusCode == http.StatusNotFound {
							log.WithFields(logrus.Fields{
								"resource": resourceType,
								"endpoint": endpoint,
							}).Debug("no resources found")

							continue
						}
					}
					log.Fatalf("failed to fetch API endpoint: %s", err)
				}

				body, err := io.ReadAll(result.Body)
				if err != nil {
					log.Fatalln(err)
				}

				value := gjson.Get(string(body), "result")
				if value.Type == gjson.Null {
					log.WithFields(logrus.Fields{
						"resource": resourceType,
						"endpoint": endpoint,
					}).Debug("no result found")

					continue
				}
				err = json.Unmarshal([]byte(value.String()), &jsonStructData)
				if err != nil {
					log.Fatalf("failed to unmarshal result: %s", err)
				}
			}
		} else {
			var identifier *cfv0.ResourceContainer
			if accountID != "" {
				identifier = cfv0.AccountIdentifier(accountID)
			} else {
				identifier = cfv0.ZoneIdentifier(zoneID)
			}

			resources := strings.Split(resourceType, ",")
			for _, resourceType := range resources {
				switch resourceType {
				case "cloudflare_access_application":
					jsonPayload, _, err := api.ListAccessApplications(context.Background(), identifier, cfv0.ListAccessApplicationsParams{})
					if err != nil {
						log.Fatal(err)
					}

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_access_group":
					jsonPayload, _, err := api.ListAccessGroups(context.Background(), identifier, cfv0.ListAccessGroupsParams{})
					if err != nil {
						log.Fatal(err)
					}

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_access_rule":
					if accountID != "" {
						jsonPayload, err := api.ListAccountAccessRules(context.Background(), accountID, cfv0.AccessRule{}, 1)
						if err != nil {
							log.Fatal(err)
						}

						m, _ := json.Marshal(jsonPayload.Result)
						err = json.Unmarshal(m, &jsonStructData)
						if err != nil {
							log.Fatal(err)
						}
					} else {
						jsonPayload, err := api.ListZoneAccessRules(context.Background(), zoneID, cfv0.AccessRule{}, 1)
						if err != nil {
							log.Fatal(err)
						}

						m, _ := json.Marshal(jsonPayload.Result)
						err = json.Unmarshal(m, &jsonStructData)
						if err != nil {
							log.Fatal(err)
						}
					}
				case "cloudflare_account_member":
					jsonPayload, _, err := api.AccountMembers(context.Background(), accountID, cfv0.PaginationOptions{})
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_argo":
					jsonPayload := []cfv0.ArgoFeatureSetting{{
						ID: fmt.Sprintf("%x", md5.Sum([]byte(time.Now().String()))),
					}}

					m, _ := json.Marshal(jsonPayload)
					err := json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_bot_management":
					botManagement, err := api.GetBotManagement(context.Background(), identifier)
					if err != nil {
						log.Fatal(err)
					}
					var jsonPayload []cfv0.BotManagement
					jsonPayload = append(jsonPayload, botManagement)

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					jsonStructData[0].(map[string]interface{})["id"] = zoneID
				case "cloudflare_byo_ip_prefix":
					jsonPayload, err := api.ListPrefixes(context.Background(), accountID)
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_certificate_pack":
					jsonPayload, err := api.ListCertificatePacks(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}

					var customerManagedCertificates []cfv0.CertificatePack
					for _, r := range jsonPayload {
						if r.Type != "universal" {
							customerManagedCertificates = append(customerManagedCertificates, r)
						}
					}
					jsonPayload = customerManagedCertificates

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_custom_pages":
					if accountID != "" {
						jsonPayload, err := api.CustomPages(context.Background(), &cfv0.CustomPageOptions{AccountID: accountID})
						if err != nil {
							log.Fatal(err)
						}

						m, _ := json.Marshal(jsonPayload)
						err = json.Unmarshal(m, &jsonStructData)
						if err != nil {
							log.Fatal(err)
						}
					} else {
						jsonPayload, err := api.CustomPages(context.Background(), &cfv0.CustomPageOptions{ZoneID: zoneID})
						if err != nil {
							log.Fatal(err)
						}

						m, _ := json.Marshal(jsonPayload)
						err = json.Unmarshal(m, &jsonStructData)
						if err != nil {
							log.Fatal(err)
						}
					}
				case "cloudflare_filter":
					jsonPayload, _, err := api.Filters(context.Background(), identifier, cfv0.FilterListParams{})
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_firewall_rule":
					jsonPayload, _, err := api.FirewallRules(context.Background(), identifier, cfv0.FirewallRuleListParams{})
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_healthcheck":
					jsonPayload, err := api.Healthchecks(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_custom_hostname":
					jsonPayload, _, err := api.CustomHostnames(context.Background(), zoneID, 1, cfv0.CustomHostname{})
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_custom_ssl":
					jsonPayload, err := api.ListSSL(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_ip_list":
					jsonPayload, err := api.ListIPLists(context.Background(), accountID)
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_load_balancer":
					jsonPayload, err := api.ListLoadBalancers(context.Background(), identifier, cfv0.ListLoadBalancerParams{})
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_load_balancer_pool":
					jsonPayload, err := api.ListLoadBalancerPools(context.Background(), identifier, cfv0.ListLoadBalancerPoolParams{})
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_load_balancer_monitor":
					jsonPayload, err := api.ListLoadBalancerMonitors(context.Background(), identifier, cfv0.ListLoadBalancerMonitorParams{})
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_logpush_job":
					jsonPayload, err := api.ListLogpushJobs(context.Background(), identifier, cfv0.ListLogpushJobsParams{})
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_origin_ca_certificate":
					jsonPayload, err := api.ListOriginCACertificates(context.Background(), cfv0.ListOriginCertificatesParams{ZoneID: zoneID})
					if err != nil {
						log.Fatal(err)
					}

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_page_rule":
					jsonPayload, err := api.ListPageRules(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_rate_limit":
					jsonPayload, err := api.ListAllRateLimits(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_record":
					jsonPayload, _, err := api.ListDNSRecords(context.Background(), identifier, cfv0.ListDNSRecordsParams{})
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_ruleset":
					jsonPayload, err := api.ListRulesets(context.Background(), identifier, cfv0.ListRulesetsParams{})
					if err != nil {
						log.Fatal(err)
					}

					// Customers can read-only Managed Rulesets, so we don't want to
					// have them try to import something they can't manage with terraform
					var nonManagedRules []cfv0.Ruleset
					for _, r := range jsonPayload {
						if r.Kind != string(cfv0.RulesetKindManaged) {
							nonManagedRules = append(nonManagedRules, r)
						}
					}

					m, _ := json.Marshal(nonManagedRules)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_spectrum_application":
					jsonPayload, err := api.SpectrumApplications(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_teams_list":
					jsonPayload, _, err := api.ListTeamsLists(context.Background(), identifier, cfv0.ListTeamListsParams{})
					if err != nil {
						log.Fatal(err)
					}

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_teams_location":
					jsonPayload, _, err := api.TeamsLocations(context.Background(), accountID)
					if err != nil {
						log.Fatal(err)
					}

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_teams_proxy_endpoint":
					jsonPayload, _, err := api.TeamsProxyEndpoints(context.Background(), accountID)
					if err != nil {
						log.Fatal(err)
					}

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_teams_rule":
					jsonPayload, err := api.TeamsRules(context.Background(), accountID)
					if err != nil {
						log.Fatal(err)
					}

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_tunnel":
					log.Debug("only requesting the first 1000 active Cloudflare Tunnels due to the service not providing correct pagination responses")
					jsonPayload, _, err := api.ListTunnels(
						context.Background(),
						cfv0.AccountIdentifier(accountID),
						cfv0.TunnelListParams{
							IsDeleted: cfv0.BoolPtr(false),
							ResultInfo: cfv0.ResultInfo{
								PerPage: 1000,
								Page:    1,
							},
						})
					if err != nil {
						log.Fatal(err)
					}

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_turnstile_widget":
					jsonPayload, _, err := api.ListTurnstileWidgets(context.Background(), identifier, cfv0.ListTurnstileWidgetParams{})
					if err != nil {
						log.Fatal(err)
					}

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
					for i := 0; i < len(jsonStructData); i++ {
						jsonStructData[i].(map[string]interface{})["id"] = jsonStructData[i].(map[string]interface{})["sitekey"]
					}
				case "cloudflare_waf_override":
					jsonPayload, err := api.ListWAFOverrides(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_waf_package":
					jsonPayload, err := api.ListWAFPackages(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_waiting_room":
					jsonPayload, err := api.ListWaitingRooms(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_workers_kv_namespace":
					jsonPayload, _, err := api.ListWorkersKVNamespaces(context.Background(), identifier, cfv0.ListWorkersKVNamespacesParams{})
					if err != nil {
						log.Fatal(err)
					}

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_worker_route":
					jsonPayload, err := api.ListWorkerRoutes(context.Background(), identifier, cfv0.ListWorkerRoutesParams{})
					if err != nil {
						log.Fatal(err)
					}

					m, _ := json.Marshal(jsonPayload.Routes)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_zone":
					jsonPayload, err := api.ListZones(context.Background())
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_zone_lockdown":
					jsonPayload, _, err := api.ListZoneLockdowns(context.Background(), identifier, cfv0.LockdownListParams{})
					if err != nil {
						log.Fatal(err)
					}

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				default:
					fmt.Fprintf(cmd.OutOrStderr(), "%q is not yet supported for state import", resourceType)
					return
				}
			}
		}

		importFile := hclwrite.NewEmptyFile()
		importBody := importFile.Body()

		for _, data := range jsonStructData {
			id := data.(map[string]interface{})["id"].(string)

			if useModernImportBlock {
				idvalue := buildRawImportAddress(resourceType, id, resourceToEndpoint[resourceType]["get"])
				imp := importBody.AppendNewBlock("import", []string{}).Body()
				imp.SetAttributeRaw("to", hclwrite.TokensForIdentifier(fmt.Sprintf("%s.%s", resourceType, fmt.Sprintf("%s_%s", terraformResourceNamePrefix, id))))
				imp.SetAttributeValue("id", cty.StringVal(idvalue))
				importFile.Body().AppendNewline()
			} else {
				fmt.Fprint(cmd.OutOrStdout(), buildTerraformImportCommand(resourceType, id, resourceToEndpoint[resourceType]["get"]))
			}
		}

		if useModernImportBlock {
			// don't format the output; there is a bug in hclwrite.Format that
			// splits incorrectly on certain characters. instead, manually
			// insert new lines on the block.
			fmt.Fprint(cmd.OutOrStdout(), string(importFile.Bytes()))
		}
	}
}

// buildTerraformImportCommand takes the resourceType and resourceID in order to
// lookup the resource type import string and then return a suitable composite
// value that is compatible with `terraform import`.
//
// Note: `endpoint` is only used on > v4. Otherwise it is ignored.
func buildTerraformImportCommand(resourceType, resourceID, endpoint string) string {
	resourceImportAddress := buildRawImportAddress(resourceType, resourceID, endpoint)
	return fmt.Sprintf("%s %s.%s_%s %s\n", terraformImportCmdPrefix, resourceType, terraformResourceNamePrefix, resourceID, resourceImportAddress)
}

// buildRawImportAddress takes the resourceType and resourceID in order to lookup
// the resource type import string and then return a suitable address.
//
// Note: `endpoint` is only used on > v4. Otherwise it is ignored.
func buildRawImportAddress(resourceType, resourceID, endpoint string) string {
	if strings.HasPrefix(providerVersionString, "5") {
		prefix := ""
		if strings.Contains(endpoint, "{account_or_zone}") {
			if accountID != "" {
				prefix = "accounts"
				endpoint = strings.Replace(endpoint, "/{account_or_zone}/{account_or_zone_id}/", "/accounts/{account_id}/", 1)
			} else {
				prefix = "zones"
				endpoint = strings.Replace(endpoint, "/{account_or_zone}/{account_or_zone_id}/", "/zones/{zone_id}/", 1)
			}
		}

		r, _ := regexp.Compile("({[a-z0-9_]*})")
		matches := r.FindAllString(endpoint, -1)

		if len(matches) == 1 {
			matches[0] = resourceID
		} else {
			matches[1] = resourceID
		}

		output := strings.Join(matches, "/")

		replacer := strings.NewReplacer(
			"{account_id}", accountID,
			"{zone_id}", zoneID,
		)

		if prefix != "" {
			output = prefix + "/" + output
		}

		return replacer.Replace(output)
	} else {
		if _, ok := resourceImportStringFormats[resourceType]; !ok {
			log.Fatalf("%s does not have an import format defined", resourceType)
		}

		var identiferType string
		var identiferValue string

		if accountID != "" {
			identiferType = "account"
			identiferValue = accountID
		} else {
			identiferType = "zone"
			identiferValue = zoneID
		}

		s := resourceImportStringFormats[resourceType]
		replacer := strings.NewReplacer(
			":identifier_type", identiferType,
			":identifier_value", identiferValue,
			":zone_id", zoneID,
			":account_id", accountID,
			":id", resourceID,
		)

		return replacer.Replace(s)
	}
}
