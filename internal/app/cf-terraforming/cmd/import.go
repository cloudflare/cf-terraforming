package cmd

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/spf13/cobra"
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
		var jsonStructData []interface{}

		var identifier *cloudflare.ResourceContainer
		if accountID != "" {
			identifier = cloudflare.AccountIdentifier(accountID)
		} else {
			identifier = cloudflare.ZoneIdentifier(zoneID)
		}

		resources := strings.Split(resourceType, ",")
		for _, resourceType := range resources {
			switch resourceType {
			case "cloudflare_access_application":
				jsonPayload, _, err := api.ListAccessApplications(context.Background(), identifier, cloudflare.ListAccessApplicationsParams{})
				if err != nil {
					log.Fatal(err)
				}

				m, _ := json.Marshal(jsonPayload)
				err = json.Unmarshal(m, &jsonStructData)
				if err != nil {
					log.Fatal(err)
				}
			case "cloudflare_access_group":
				jsonPayload, _, err := api.ListAccessGroups(context.Background(), identifier, cloudflare.ListAccessGroupsParams{})
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
					jsonPayload, err := api.ListAccountAccessRules(context.Background(), accountID, cloudflare.AccessRule{}, 1)
					if err != nil {
						log.Fatal(err)
					}

					m, _ := json.Marshal(jsonPayload.Result)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				} else {
					jsonPayload, err := api.ListZoneAccessRules(context.Background(), zoneID, cloudflare.AccessRule{}, 1)
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
				jsonPayload, _, err := api.AccountMembers(context.Background(), accountID, cloudflare.PaginationOptions{})
				if err != nil {
					log.Fatal(err)
				}
				m, _ := json.Marshal(jsonPayload)
				err = json.Unmarshal(m, &jsonStructData)
				if err != nil {
					log.Fatal(err)
				}
			case "cloudflare_argo":
				jsonPayload := []cloudflare.ArgoFeatureSetting{{
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
				var jsonPayload []cloudflare.BotManagement
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

				var customerManagedCertificates []cloudflare.CertificatePack
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
					jsonPayload, err := api.CustomPages(context.Background(), &cloudflare.CustomPageOptions{AccountID: accountID})
					if err != nil {
						log.Fatal(err)
					}

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				} else {
					jsonPayload, err := api.CustomPages(context.Background(), &cloudflare.CustomPageOptions{ZoneID: zoneID})
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
				jsonPayload, _, err := api.Filters(context.Background(), identifier, cloudflare.FilterListParams{})
				if err != nil {
					log.Fatal(err)
				}
				m, _ := json.Marshal(jsonPayload)
				err = json.Unmarshal(m, &jsonStructData)
				if err != nil {
					log.Fatal(err)
				}
			case "cloudflare_firewall_rule":
				jsonPayload, _, err := api.FirewallRules(context.Background(), identifier, cloudflare.FirewallRuleListParams{})
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
				jsonPayload, _, err := api.CustomHostnames(context.Background(), zoneID, 1, cloudflare.CustomHostname{})
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
				jsonPayload, err := api.ListLoadBalancers(context.Background(), identifier, cloudflare.ListLoadBalancerParams{})
				if err != nil {
					log.Fatal(err)
				}
				m, _ := json.Marshal(jsonPayload)
				err = json.Unmarshal(m, &jsonStructData)
				if err != nil {
					log.Fatal(err)
				}
			case "cloudflare_load_balancer_pool":
				jsonPayload, err := api.ListLoadBalancerPools(context.Background(), identifier, cloudflare.ListLoadBalancerPoolParams{})
				if err != nil {
					log.Fatal(err)
				}
				m, _ := json.Marshal(jsonPayload)
				err = json.Unmarshal(m, &jsonStructData)
				if err != nil {
					log.Fatal(err)
				}
			case "cloudflare_load_balancer_monitor":
				jsonPayload, err := api.ListLoadBalancerMonitors(context.Background(), identifier, cloudflare.ListLoadBalancerMonitorParams{})
				if err != nil {
					log.Fatal(err)
				}
				m, _ := json.Marshal(jsonPayload)
				err = json.Unmarshal(m, &jsonStructData)
				if err != nil {
					log.Fatal(err)
				}
			case "cloudflare_logpush_job":
				jsonPayload, err := api.ListLogpushJobs(context.Background(), identifier, cloudflare.ListLogpushJobsParams{})
				if err != nil {
					log.Fatal(err)
				}
				m, _ := json.Marshal(jsonPayload)
				err = json.Unmarshal(m, &jsonStructData)
				if err != nil {
					log.Fatal(err)
				}
			case "cloudflare_origin_ca_certificate":
				jsonPayload, err := api.ListOriginCACertificates(context.Background(), cloudflare.ListOriginCertificatesParams{ZoneID: zoneID})
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
				jsonPayload, _, err := api.ListDNSRecords(context.Background(), identifier, cloudflare.ListDNSRecordsParams{})
				if err != nil {
					log.Fatal(err)
				}
				m, _ := json.Marshal(jsonPayload)
				err = json.Unmarshal(m, &jsonStructData)
				if err != nil {
					log.Fatal(err)
				}
			case "cloudflare_ruleset":
				jsonPayload, err := api.ListRulesets(context.Background(), identifier, cloudflare.ListRulesetsParams{})
				if err != nil {
					log.Fatal(err)
				}

				// Customers can read-only Managed Rulesets, so we don't want to
				// have them try to import something they can't manage with terraform
				var nonManagedRules []cloudflare.Ruleset
				for _, r := range jsonPayload {
					if r.Kind != string(cloudflare.RulesetKindManaged) {
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
				jsonPayload, _, err := api.ListTeamsLists(context.Background(), identifier, cloudflare.ListTeamListsParams{})
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
					cloudflare.AccountIdentifier(accountID),
					cloudflare.TunnelListParams{
						IsDeleted: cloudflare.BoolPtr(false),
						ResultInfo: cloudflare.ResultInfo{
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
				jsonPayload, _, err := api.ListTurnstileWidgets(context.Background(), identifier, cloudflare.ListTurnstileWidgetParams{})
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
				jsonPayload, _, err := api.ListWorkersKVNamespaces(context.Background(), identifier, cloudflare.ListWorkersKVNamespacesParams{})
				if err != nil {
					log.Fatal(err)
				}

				m, _ := json.Marshal(jsonPayload)
				err = json.Unmarshal(m, &jsonStructData)
				if err != nil {
					log.Fatal(err)
				}
			case "cloudflare_worker_route":
				jsonPayload, err := api.ListWorkerRoutes(context.Background(), identifier, cloudflare.ListWorkerRoutesParams{})
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
				jsonPayload, _, err := api.ListZoneLockdowns(context.Background(), identifier, cloudflare.LockdownListParams{})
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

			importFile := hclwrite.NewEmptyFile()
			importBody := importFile.Body()

			for _, data := range jsonStructData {
				id := data.(map[string]interface{})["id"].(string)

				if useModernImportBlock {
					idvalue := buildRawImportAddress(resourceType, id)
					imp := importBody.AppendNewBlock("import", []string{}).Body()
					imp.SetAttributeRaw("to", hclwrite.TokensForIdentifier(fmt.Sprintf("%s.%s", resourceType, fmt.Sprintf("%s_%s", terraformResourceNamePrefix, id))))
					imp.SetAttributeValue("id", cty.StringVal(idvalue))
					importFile.Body().AppendNewline()
				} else {
					fmt.Fprint(cmd.OutOrStdout(), buildTerraformImportCommand(resourceType, id))
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
}

// buildTerraformImportCommand takes the resourceType and resourceID in order to
// lookup the resource type import string and then return a suitable composite
// value that is compatible with `terraform import`.
func buildTerraformImportCommand(resourceType, resourceID string) string {
	resourceImportAddress := buildRawImportAddress(resourceType, resourceID)
	return fmt.Sprintf("%s %s.%s_%s %s\n", terraformImportCmdPrefix, resourceType, terraformResourceNamePrefix, resourceID, resourceImportAddress)
}

// buildRawImportAddress takes the resourceType and resourceID in order to lookup
// the resource type import string and then return a suitable address.
func buildRawImportAddress(resourceType, resourceID string) string {
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
