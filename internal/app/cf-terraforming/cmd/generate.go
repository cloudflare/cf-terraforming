package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	cfv0 "github.com/cloudflare/cloudflare-go"
	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"github.com/zclconf/go-cty/cty"
)

var resourceType string

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:    "generate",
	Short:  "Fetch resources from the Cloudflare API and generate the respective Terraform stanzas",
	Run:    generateResources(),
	PreRun: sharedPreRun,
}

func generateResources() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if resourceType == "" {
			log.Fatal("you must define a resource type to generate")
		}

		zoneID = viper.GetString("zone")
		accountID = viper.GetString("account")
		workingDir := viper.GetString("terraform-install-path")
		execPath := viper.GetString("terraform-binary-path")
		providerRegistryHostname := viper.GetString("provider-registry-hostname")

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
		providerVersionString := providerVersion[providerRegistryHostname+"/cloudflare/cloudflare"].String()
		log.Debugf("detected provider version: %s", providerVersionString)

		log.Debug("reading Terraform schema for Cloudflare provider")
		ps, err := tf.ProvidersSchema(context.Background())
		if err != nil {
			log.Fatal("failed to read provider schema", err)
		}

		s := ps.Schemas[providerRegistryHostname+"/cloudflare/cloudflare"]
		if s == nil {
			log.Fatal("failed to detect provider installation")
		}

		resources := strings.Split(resourceType, ",")
		for _, resourceType := range resources {
			r := s.ResourceSchemas[resourceType]
			log.Debugf("beginning to read and build %q resources", resourceType)

			// Initialise `resourceCount` outside of the switch for supported resources
			// to allow it to be referenced further down in the loop that outputs the
			// newly generated resources.
			resourceCount := 0
			var jsonStructData []interface{}

			if strings.HasPrefix(providerVersionString, "5") {
				if resourceToEndpoint[resourceType] == "" {
					log.Debugf("did not find API endpoint for %q. skipping...", resourceType)
					continue
				}

				var result *http.Response

				endpoint := resourceToEndpoint[resourceType]

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

				client := cloudflare.NewClient()

				err := client.Get(context.Background(), endpoint, nil, &result)
				if err != nil {
					log.Fatalf("failed to fetch API endpoint: %s", err)
				}

				body, err := io.ReadAll(result.Body)
				if err != nil {
					log.Fatalln(err)
				}

				value := gjson.Get(string(body), "result")
				err = json.Unmarshal([]byte(value.String()), &jsonStructData)
				if err != nil {
					log.Fatalf("failed to unmarshal result: %s", err)
				}

				resourceCount = len(jsonStructData)
			} else {
				var identifier *cfv0.ResourceContainer
				if accountID != "" {
					identifier = cfv0.AccountIdentifier(accountID)
				} else {
					identifier = cfv0.ZoneIdentifier(zoneID)
				}

				switch resourceType {
				case "cloudflare_access_application":
					jsonPayload, _, err := api.ListAccessApplications(context.Background(), identifier, cfv0.ListAccessApplicationsParams{})
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
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

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_access_identity_provider":
					jsonPayload, _, err := api.ListAccessIdentityProviders(context.Background(), identifier, cfv0.ListAccessIdentityProvidersParams{})
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_access_service_token":
					jsonPayload, _, err := api.ListAccessServiceTokens(context.Background(), identifier, cfv0.ListAccessServiceTokensParams{})
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_access_mutual_tls_certificate":
					jsonPayload, _, err := api.ListAccessMutualTLSCertificates(context.Background(), identifier, cfv0.ListAccessMutualTLSCertificatesParams{})
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
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

						resourceCount = len(jsonPayload.Result)
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

						resourceCount = len(jsonPayload.Result)
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

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					// remap email and role_ids into the right structure.
					for i := 0; i < resourceCount; i++ {
						jsonStructData[i].(map[string]interface{})["email_address"] = jsonStructData[i].(map[string]interface{})["user"].(map[string]interface{})["email"]
						roleIDs := []string{}
						for _, role := range jsonStructData[i].(map[string]interface{})["roles"].([]interface{}) {
							roleIDs = append(roleIDs, role.(map[string]interface{})["id"].(string))
						}
						jsonStructData[i].(map[string]interface{})["role_ids"] = roleIDs
					}
				case "cloudflare_argo":
					jsonPayload := []cfv0.ArgoFeatureSetting{}

					argoSmartRouting, err := api.ArgoSmartRouting(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}
					jsonPayload = append(jsonPayload, argoSmartRouting)

					argoTieredCaching, err := api.ArgoTieredCaching(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}
					jsonPayload = append(jsonPayload, argoTieredCaching)

					resourceCount = 1

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					for i, b := range jsonStructData {
						key := b.(map[string]interface{})["id"].(string)
						jsonStructData[0].(map[string]interface{})[key] = jsonStructData[i].(map[string]interface{})["value"]
					}
				case "cloudflare_api_shield":
					jsonPayload := []cfv0.APIShield{}
					apiShieldConfig, _, err := api.GetAPIShieldConfiguration(context.Background(), identifier)
					if err != nil {
						log.Fatal(err)
					}
					// the response can contain an empty APIShield struct. Verify we have data before we attempt to do anything
					jsonPayload = append(jsonPayload, apiShieldConfig)

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					// this is only every a 1:1 so we can just verify if the 0th element has they key we expect
					jsonStructData[0].(map[string]interface{})["id"] = zoneID

					if jsonStructData[0].(map[string]interface{})["auth_id_characteristics"] == nil {
						// force a no resources return by setting resourceCount to 0
						resourceCount = 0
					}
				case "cloudflare_user_agent_blocking_rule":
					page := 1
					var jsonPayload []cfv0.UserAgentRule
					for {
						res, err := api.ListUserAgentRules(context.Background(), zoneID, page)
						if err != nil {
							log.Fatal(err)
						}

						jsonPayload = append(jsonPayload, res.Result...)
						res.ResultInfo = res.ResultInfo.Next()

						if res.ResultInfo.Done() {
							break
						}
						page = page + 1
					}

					resourceCount = len(jsonPayload)
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

					resourceCount = 1
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

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					// remap ID to prefix_id and advertised to advertisement on the JSON payloads.
					for i := 0; i < resourceCount; i++ {
						jsonStructData[i].(map[string]interface{})["prefix_id"] = jsonStructData[i].(map[string]interface{})["id"]

						if jsonStructData[i].(map[string]interface{})["advertised"].(bool) {
							jsonStructData[i].(map[string]interface{})["advertisement"] = "on"
						} else {
							jsonStructData[i].(map[string]interface{})["advertisement"] = "off"
						}
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

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_custom_pages":
					if accountID != "" {
						acc := cfv0.CustomPageOptions{AccountID: accountID}
						jsonPayload, err := api.CustomPages(context.Background(), &acc)
						if err != nil {
							log.Fatal(err)
						}

						resourceCount = len(jsonPayload)
						m, _ := json.Marshal(jsonPayload)
						err = json.Unmarshal(m, &jsonStructData)
						if err != nil {
							log.Fatal(err)
						}
					} else {
						zo := cfv0.CustomPageOptions{ZoneID: zoneID}
						jsonPayload, err := api.CustomPages(context.Background(), &zo)
						if err != nil {
							log.Fatal(err)
						}

						resourceCount = len(jsonPayload)
						m, _ := json.Marshal(jsonPayload)
						err = json.Unmarshal(m, &jsonStructData)
						if err != nil {
							log.Fatal(err)
						}
					}

					var newJsonStructData []interface{}
					// remap ID to the "type" field
					for i := 0; i < resourceCount; i++ {
						jsonStructData[i].(map[string]interface{})["type"] = jsonStructData[i].(map[string]interface{})["id"]
						// we only want repsonses that have 'url'
						if jsonStructData[i].(map[string]interface{})["url"] != nil {
							newJsonStructData = append(newJsonStructData, jsonStructData[i])
						}
					}
					jsonStructData = newJsonStructData
					resourceCount = len(jsonStructData)

				case "cloudflare_custom_hostname_fallback_origin":
					var jsonPayload []cfv0.CustomHostnameFallbackOrigin
					apiCall, err := api.CustomHostnameFallbackOrigin(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}

					if apiCall.Origin != "" {
						resourceCount = 1
						jsonPayload = append(jsonPayload, apiCall)
					}

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					for i := 0; i < resourceCount; i++ {
						jsonStructData[i].(map[string]interface{})["id"] = sanitiseTerraformResourceName(jsonStructData[i].(map[string]interface{})["origin"].(string))
						jsonStructData[i].(map[string]interface{})["status"] = nil
					}
				case "cloudflare_filter":
					jsonPayload, _, err := api.Filters(context.Background(), identifier, cfv0.FilterListParams{})
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
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

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					// remap Filter.ID to `filter_id` on the JSON payloads.
					for i := 0; i < resourceCount; i++ {
						jsonStructData[i].(map[string]interface{})["filter_id"] = jsonStructData[i].(map[string]interface{})["filter"].(map[string]interface{})["id"]
					}
				case "cloudflare_custom_hostname":
					jsonPayload, _, err := api.CustomHostnames(context.Background(), zoneID, 1, cfv0.CustomHostname{})
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					for i := 0; i < resourceCount; i++ {
						jsonStructData[i].(map[string]interface{})["ssl"].(map[string]interface{})["validation_errors"] = nil
					}
				case "cloudflare_custom_ssl":
					jsonPayload, err := api.ListSSL(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
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

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_list":
					jsonPayload, err := api.ListLists(context.Background(), identifier, cfv0.ListListsParams{})
					if err != nil {
						log.Fatal(err)
					}

					m, err := json.Marshal(jsonPayload)
					if err != nil {
						log.Fatal(err)
					}

					if err = json.Unmarshal(m, &jsonStructData); err != nil {
						log.Fatal(err)
					}
					resourceCount = len(jsonPayload)

					for i := 0; i < resourceCount; i++ {
						listID := jsonPayload[i].ID
						kind := jsonPayload[i].Kind

						listItems, err := api.ListListItems(context.Background(), identifier, cfv0.ListListItemsParams{ID: listID})
						if err != nil {
							log.Fatal(err)
						}
						items := make([]interface{}, 0)

						for _, listItem := range listItems {
							if kind == "" {
								continue
							}

							value := map[string]interface{}{}
							switch kind {
							case "ip":
								if listItem.IP == nil {
									continue
								}
								value["ip"] = *listItem.IP
							case "asn":
								if listItem.ASN == nil {
									continue
								}
								value["asn"] = int(*listItem.ASN)
							case "hostname":
								if listItem.Hostname == nil {
									continue
								}
								value["hostname"] = map[string]interface{}{
									"url_hostname": listItem.Hostname.UrlHostname,
								}
							case "redirect":
								if listItem.Redirect == nil {
									continue
								}
								redirect := map[string]interface{}{
									"source_url": listItem.Redirect.SourceUrl,
									"target_url": listItem.Redirect.TargetUrl,
								}
								if listItem.Redirect.IncludeSubdomains != nil {
									redirect["include_subdomains"] = boolToEnabledOrDisabled(*listItem.Redirect.IncludeSubdomains)
								}
								if listItem.Redirect.SubpathMatching != nil {
									redirect["subpath_matching"] = boolToEnabledOrDisabled(*listItem.Redirect.SubpathMatching)
								}
								if listItem.Redirect.StatusCode != nil {
									redirect["status_code"] = *listItem.Redirect.StatusCode
								}
								if listItem.Redirect.PreserveQueryString != nil {
									redirect["preserve_query_string"] = boolToEnabledOrDisabled(*listItem.Redirect.PreserveQueryString)
								}
								if listItem.Redirect.PreservePathSuffix != nil {
									redirect["preserve_path_suffix"] = boolToEnabledOrDisabled(*listItem.Redirect.PreservePathSuffix)
								}
								value["redirect"] = redirect
							}
							items = append(items, map[string]interface{}{
								"comment": listItem.Comment,
								"value":   value,
							})
						}
						jsonStructData[i].(map[string]interface{})["item"] = items
					}
				case "cloudflare_load_balancer":
					jsonPayload, err := api.ListLoadBalancers(context.Background(), identifier, cfv0.ListLoadBalancerParams{})
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					for i := 0; i < resourceCount; i++ {
						jsonStructData[i].(map[string]interface{})["default_pool_ids"] = jsonStructData[i].(map[string]interface{})["default_pools"]
						jsonStructData[i].(map[string]interface{})["fallback_pool_id"] = jsonStructData[i].(map[string]interface{})["fallback_pool"]

						if jsonStructData[i].(map[string]interface{})["country_pools"] != nil {
							original := jsonStructData[i].(map[string]interface{})["country_pools"]
							jsonStructData[i].(map[string]interface{})["country_pools"] = []interface{}{}

							for country, popIDs := range original.(map[string]interface{}) {
								jsonStructData[i].(map[string]interface{})["country_pools"] = append(jsonStructData[i].(map[string]interface{})["country_pools"].([]interface{}), map[string]interface{}{"country": country, "pool_ids": popIDs})
							}
						}

						if jsonStructData[i].(map[string]interface{})["region_pools"] != nil {
							original := jsonStructData[i].(map[string]interface{})["region_pools"]
							jsonStructData[i].(map[string]interface{})["region_pools"] = []interface{}{}

							for region, popIDs := range original.(map[string]interface{}) {
								jsonStructData[i].(map[string]interface{})["region_pools"] = append(jsonStructData[i].(map[string]interface{})["region_pools"].([]interface{}), map[string]interface{}{"region": region, "pool_ids": popIDs})
							}
						}

						if jsonStructData[i].(map[string]interface{})["pop_pools"] != nil {
							original := jsonStructData[i].(map[string]interface{})["pop_pools"]
							jsonStructData[i].(map[string]interface{})["pop_pools"] = []interface{}{}

							for pop, popIDs := range original.(map[string]interface{}) {
								jsonStructData[i].(map[string]interface{})["pop_pools"] = append(jsonStructData[i].(map[string]interface{})["pop_pools"].([]interface{}), map[string]interface{}{"pop": pop, "pool_ids": popIDs})
							}
						}
					}

				case "cloudflare_load_balancer_pool":
					jsonPayload, err := api.ListLoadBalancerPools(context.Background(), identifier, cfv0.ListLoadBalancerPoolParams{})
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					for i := 0; i < resourceCount; i++ {
						for originCounter := range jsonStructData[i].(map[string]interface{})["origins"].([]interface{}) {
							if jsonStructData[i].(map[string]interface{})["origins"].([]interface{})[originCounter].(map[string]interface{})["header"] != nil {
								jsonStructData[i].(map[string]interface{})["origins"].([]interface{})[originCounter].(map[string]interface{})["header"].(map[string]interface{})["header"] = "Host"
								jsonStructData[i].(map[string]interface{})["origins"].([]interface{})[originCounter].(map[string]interface{})["header"].(map[string]interface{})["values"] = jsonStructData[i].(map[string]interface{})["origins"].([]interface{})[originCounter].(map[string]interface{})["header"].(map[string]interface{})["Host"]
							}
						}
					}
				case "cloudflare_load_balancer_monitor":
					jsonPayload, err := api.ListLoadBalancerMonitors(context.Background(), identifier, cfv0.ListLoadBalancerMonitorParams{})
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
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

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					for i := 0; i < resourceCount; i++ {
						// Workaround for LogpushJob.Filter being empty with a custom
						// marshaler and returning `{"where":{}}` as the "empty" value.
						if jsonStructData[i].(map[string]interface{})["filter"] == `{"where":{}}` {
							jsonStructData[i].(map[string]interface{})["filter"] = nil
						}
					}
				case "cloudflare_managed_headers":
					// only grab the enabled headers
					jsonPayload, err := api.ListZoneManagedHeaders(context.Background(), cfv0.ResourceIdentifier(zoneID), cfv0.ListManagedHeadersParams{Status: "enabled"})
					if err != nil {
						log.Fatal(err)
					}

					var managedHeaders []cfv0.ManagedHeaders
					managedHeaders = append(managedHeaders, jsonPayload)

					resourceCount = len(managedHeaders)
					m, _ := json.Marshal(managedHeaders)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					for i := 0; i < resourceCount; i++ {
						jsonStructData[i].(map[string]interface{})["id"] = zoneID
					}
				case "cloudflare_origin_ca_certificate":
					jsonPayload, err := api.ListOriginCACertificates(context.Background(), cfv0.ListOriginCertificatesParams{ZoneID: zoneID})
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
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

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					for i := 0; i < resourceCount; i++ {
						jsonStructData[i].(map[string]interface{})["target"] = jsonStructData[i].(map[string]interface{})["targets"].([]interface{})[0].(map[string]interface{})["constraint"].(map[string]interface{})["value"]
						jsonStructData[i].(map[string]interface{})["actions"] = flattenAttrMap(jsonStructData[i].(map[string]interface{})["actions"].([]interface{}))

						// Have to remap the cache_ttl_by_status to conform to Terraform's more human-friendly structure.
						if cache, ok := jsonStructData[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_ttl_by_status"].(map[string]interface{}); ok {
							cache_ttl_by_status := []map[string]interface{}{}

							for codes, ttl := range cache {
								if ttl == "no-cache" {
									ttl = 0
								} else if ttl == "no-store" {
									ttl = -1
								}
								elem := map[string]interface{}{
									"codes": codes,
									"ttl":   ttl,
								}

								cache_ttl_by_status = append(cache_ttl_by_status, elem)
							}

							sort.SliceStable(cache_ttl_by_status, func(i int, j int) bool {
								return cache_ttl_by_status[i]["codes"].(string) < cache_ttl_by_status[j]["codes"].(string)
							})

							jsonStructData[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_ttl_by_status"] = cache_ttl_by_status
						}

						// Remap cache_key_fields.query_string.include & .exclude wildcards (not in an array) to the appropriate "ignore" field value in Terraform.
						if c, ok := jsonStructData[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_key_fields"].(map[string]interface{}); ok {
							if s, sok := c["query_string"].(map[string]interface{})["include"].(string); sok && s == "*" {
								jsonStructData[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_key_fields"].(map[string]interface{})["query_string"].(map[string]interface{})["include"] = nil
								jsonStructData[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_key_fields"].(map[string]interface{})["query_string"].(map[string]interface{})["ignore"] = false
							}
							if s, sok := c["query_string"].(map[string]interface{})["exclude"].(string); sok && s == "*" {
								jsonStructData[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_key_fields"].(map[string]interface{})["query_string"].(map[string]interface{})["exclude"] = nil
								jsonStructData[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_key_fields"].(map[string]interface{})["query_string"].(map[string]interface{})["ignore"] = true
							}
						}
					}
				case "cloudflare_rate_limit":
					jsonPayload, err := api.ListAllRateLimits(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					for i := 0; i < resourceCount; i++ {
						var bypassItems []string

						// Remap match.request.url to match.request.url_pattern
						jsonStructData[i].(map[string]interface{})["match"].(map[string]interface{})["request"].(map[string]interface{})["url_pattern"] = jsonStructData[i].(map[string]interface{})["match"].(map[string]interface{})["request"].(map[string]interface{})["url"]

						// Remap bypass to bypass_url_patterns
						if jsonStructData[i].(map[string]interface{})["bypass"] != nil {
							for _, item := range jsonStructData[i].(map[string]interface{})["bypass"].([]interface{}) {
								bypassItems = append(bypassItems, item.(map[string]interface{})["value"].(string))
							}
							jsonStructData[i].(map[string]interface{})["bypass_url_patterns"] = bypassItems
						}

						// Remap match.response.status to match.response.statuses
						jsonStructData[i].(map[string]interface{})["match"].(map[string]interface{})["response"].(map[string]interface{})["statuses"] = jsonStructData[i].(map[string]interface{})["match"].(map[string]interface{})["response"].(map[string]interface{})["status"]
					}

				case "cloudflare_record":
					jsonPayload, _, err := api.ListDNSRecords(context.Background(), identifier, cfv0.ListDNSRecordsParams{})
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					zone, _ := api.ZoneDetails(context.Background(), identifier.Identifier)

					for i := 0; i < resourceCount; i++ {
						// Drop the proxiable values as they are not usable
						jsonStructData[i].(map[string]interface{})["proxiable"] = nil
						jsonStructData[i].(map[string]interface{})["value"] = nil

						if jsonStructData[i].(map[string]interface{})["name"].(string) != zone.Name {
							jsonStructData[i].(map[string]interface{})["name"] = strings.ReplaceAll(jsonStructData[i].(map[string]interface{})["name"].(string), "."+zone.Name, "")
						}
					}
				case "cloudflare_ruleset":
					jsonPayload, err := api.ListRulesets(context.Background(), identifier, cfv0.ListRulesetsParams{})
					if err != nil {
						log.Fatal(err)
					}

					var nonManagedRules []cfv0.Ruleset

					// A little annoying but makes more sense doing it this way. Only append
					// the non-managed rules to the usable nonManagedRules variable instead
					// of attempting to delete from an existing slice and just reassign.
					for _, r := range jsonPayload {
						if r.Kind != string(cfv0.RulesetKindManaged) {
							nonManagedRules = append(nonManagedRules, r)
						}
					}
					jsonPayload = nonManagedRules
					ruleHeaders := map[string][]map[string]interface{}{}
					for i, rule := range nonManagedRules {
						ruleset, _ := api.GetRuleset(context.Background(), identifier, rule.ID)
						jsonPayload[i].Rules = ruleset.Rules

						if ruleset.Rules != nil {
							for _, rule := range ruleset.Rules {
								if rule.ActionParameters != nil && rule.ActionParameters.Headers != nil {
									// Sort the headers to have deterministic config output
									keys := make([]string, 0, len(rule.ActionParameters.Headers))
									for k := range rule.ActionParameters.Headers {
										keys = append(keys, k)
									}
									sort.Strings(keys)

									// The structure of the API response for headers differs from the
									// structure terraform requires. So we collect all the headers
									// indexed by rule.ID to massage the jsonStructData later
									for _, headerName := range keys {
										header := map[string]interface{}{
											"name":       headerName,
											"operation":  rule.ActionParameters.Headers[headerName].Operation,
											"expression": rule.ActionParameters.Headers[headerName].Expression,
											"value":      rule.ActionParameters.Headers[headerName].Value,
										}
										ruleHeaders[rule.ID] = append(ruleHeaders[rule.ID], header)
									}
								}
							}
						}
					}

					sort.Slice(jsonPayload, func(i, j int) bool {
						return jsonPayload[i].Phase < jsonPayload[j].Phase
					})

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					// Make the rules have the correct header structure
					for i, ruleset := range jsonStructData {
						if ruleset.(map[string]interface{})["rules"] != nil {
							for j, rule := range ruleset.(map[string]interface{})["rules"].([]interface{}) {
								ID := rule.(map[string]interface{})["id"]
								if ID != nil {
									headers, exists := ruleHeaders[ID.(string)]
									if exists {
										jsonStructData[i].(map[string]interface{})["rules"].([]interface{})[j].(map[string]interface{})["action_parameters"].(map[string]interface{})["headers"] = headers
									}
								}
							}
						}
					}

					// log custom fields specific transformation fields
					logCustomFieldsTransform := []string{"cookie_fields", "request_fields", "response_fields"}

					for i := 0; i < resourceCount; i++ {
						rules := jsonStructData[i].(map[string]interface{})["rules"]
						if rules != nil {
							for ruleCounter := range rules.([]interface{}) {
								// should the `ref` be the default `id`, don't output it
								// as we don't need to track a computed default.
								id := rules.([]interface{})[ruleCounter].(map[string]interface{})["id"]
								ref := rules.([]interface{})[ruleCounter].(map[string]interface{})["ref"]
								if id == ref {
									rules.([]interface{})[ruleCounter].(map[string]interface{})["ref"] = nil
								}

								actionParams := rules.([]interface{})[ruleCounter].(map[string]interface{})["action_parameters"]
								if actionParams != nil {
									// check for log custom fields that need to be transformed
									for _, logCustomFields := range logCustomFieldsTransform {
										// check if the field exists and make sure it has at least one element
										if actionParams.(map[string]interface{})[logCustomFields] != nil && len(actionParams.(map[string]interface{})[logCustomFields].([]interface{})) > 0 {
											// Create a new list to store the data in.
											var newLogCustomFields []interface{}
											// iterate over each of the keys and add them to a generic list
											for logCustomFieldsCounter := range actionParams.(map[string]interface{})[logCustomFields].([]interface{}) {
												newLogCustomFields = append(newLogCustomFields, actionParams.(map[string]interface{})[logCustomFields].([]interface{})[logCustomFieldsCounter].(map[string]interface{})["name"])
											}
											actionParams.(map[string]interface{})[logCustomFields] = newLogCustomFields
										}
									}

									// check if our ruleset is of action 'skip'
									if rules.([]interface{})[ruleCounter].(map[string]interface{})["action"] == "skip" {
										for rule := range actionParams.(map[string]interface{}) {
											// "rules" is the only map[string][]string we need to remap. The others are all []string and are handled naturally.
											if rule == "rules" {
												for key, value := range actionParams.(map[string]interface{})[rule].(map[string]interface{}) {
													var rulesList []string
													for _, val := range value.([]interface{}) {
														rulesList = append(rulesList, val.(string))
													}
													actionParams.(map[string]interface{})[rule].(map[string]interface{})[key] = strings.Join(rulesList, ",")
												}
											}
										}
									}

									// Cache Rules transformation
									if jsonStructData[i].(map[string]interface{})["phase"] == "http_request_cache_settings" {
										if ck, ok := rules.([]interface{})[ruleCounter].(map[string]interface{})["action_parameters"].(map[string]interface{})["cache_key"]; ok {
											if c, cok := ck.(map[string]interface{})["custom_key"]; cok {
												if qs, qok := c.(map[string]interface{})["query_string"]; qok {
													if s, sok := qs.(map[string]interface{})["include"]; sok && s == "*" {
														rules.([]interface{})[ruleCounter].(map[string]interface{})["action_parameters"].(map[string]interface{})["cache_key"].(map[string]interface{})["custom_key"].(map[string]interface{})["query_string"].(map[string]interface{})["include"] = []interface{}{"*"}
													}
													if s, sok := qs.(map[string]interface{})["exclude"]; sok && s == "*" {
														rules.([]interface{})[ruleCounter].(map[string]interface{})["action_parameters"].(map[string]interface{})["cache_key"].(map[string]interface{})["custom_key"].(map[string]interface{})["query_string"].(map[string]interface{})["exclude"] = []interface{}{"*"}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				case "cloudflare_spectrum_application":
					jsonPayload, err := api.SpectrumApplications(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
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
					// get items for the lists and add it the specific list struct
					for i, TeamsList := range jsonPayload {
						items_struct, _, err := api.ListTeamsListItems(
							context.Background(),
							identifier,
							cfv0.ListTeamsListItemsParams{ListID: TeamsList.ID})
						if err != nil {
							log.Fatal(err)
						}
						TeamsList.Items = append(TeamsList.Items, items_struct...)
						jsonPayload[i] = TeamsList
					}
					m, err := json.Marshal(jsonPayload)
					if err != nil {
						log.Fatal(err)
					}
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
					resourceCount = len(jsonPayload)

					// converting the items to value field and not the otherway around
					for i := 0; i < resourceCount; i++ {
						if jsonStructData[i].(map[string]interface{})["items"] != nil && len(jsonStructData[i].(map[string]interface{})["items"].([]interface{})) > 0 {
							// new interface for storing data
							var newItems []interface{}
							for _, item := range jsonStructData[i].(map[string]interface{})["items"].([]interface{}) {
								newItems = append(newItems, item.(map[string]interface{})["value"])
							}
							jsonStructData[i].(map[string]interface{})["items"] = newItems
						}
					}
				case "cloudflare_teams_location":
					jsonPayload, _, err := api.TeamsLocations(context.Background(), accountID)
					if err != nil {
						log.Fatal(err)
					}
					resourceCount = len(jsonPayload)
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
					resourceCount = len(jsonPayload)
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
					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
					// check for empty descriptions
					for i := 0; i < resourceCount; i++ {
						if jsonStructData[i].(map[string]interface{})["description"] == "" {
							jsonStructData[i].(map[string]interface{})["description"] = "default"
						}
					}
				case "cloudflare_tunnel":
					log.Debug("only requesting the first 1000 active Cloudflare Tunnels due to the service not providing correct pagination responses")
					jsonPayload, _, err := api.ListTunnels(
						context.Background(),
						identifier,
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

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					for i := 0; i < resourceCount; i++ {
						secret, err := api.GetTunnelToken(
							context.Background(),
							identifier,
							jsonStructData[i].(map[string]interface{})["id"].(string),
						)
						if err != nil {
							log.Fatal(err)
						}
						jsonStructData[i].(map[string]interface{})["secret"] = secret
						jsonStructData[i].(map[string]interface{})["account_id"] = accountID

						jsonStructData[i].(map[string]interface{})["connections"] = nil
					}
				case "cloudflare_turnstile_widget":
					jsonPayload, _, err := api.ListTurnstileWidgets(context.Background(), identifier, cfv0.ListTurnstileWidgetParams{})
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					for i := 0; i < resourceCount; i++ {
						jsonStructData[i].(map[string]interface{})["id"] = jsonStructData[i].(map[string]interface{})["sitekey"]

						// We always want to emit a list of domains, even if it is empty.
						// The empty list is used to enable the "Allow on any hostname" feature, it is *not* a default value.
						if jsonStructData[i].(map[string]interface{})["domains"] == nil {
							jsonStructData[i].(map[string]interface{})["domains"] = []string{}
						}
					}
				case "cloudflare_url_normalization_settings":
					jsonPayload, err := api.URLNormalizationSettings(context.Background(), &cfv0.ResourceContainer{Identifier: zoneID, Level: cfv0.ZoneRouteLevel})
					if err != nil {
						log.Fatal(err)
					}
					var newJsonPayload []interface{}
					newJsonPayload = append(newJsonPayload, jsonPayload)
					resourceCount = len(newJsonPayload)
					m, _ := json.Marshal(newJsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					// this is only every a 1:1 so we can just verify if the 0th element has they key we expect
					jsonStructData[0].(map[string]interface{})["id"] = zoneID
				case "cloudflare_waiting_room":
					jsonPayload, err := api.ListWaitingRooms(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}
					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					for i := 0; i < resourceCount; i++ {
						if jsonStructData[i].(map[string]interface{})["queueing_status_code"].(float64) == 0 {
							jsonStructData[i].(map[string]interface{})["queueing_status_code"] = nil
						}
					}
				case "cloudflare_waiting_room_event":
					waitingRooms, err := api.ListWaitingRooms(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}
					for i := 0; i < len(waitingRooms); i++ {
						roomEvents, err := api.ListWaitingRoomEvents(context.Background(), zoneID, waitingRooms[i].ID)
						if err != nil {
							log.Fatal(err)
						}
						m, err := json.Marshal(roomEvents)
						if err != nil {
							log.Fatal(err)
						}
						jsonRoomEvents := []interface{}{}
						err = json.Unmarshal(m, &jsonRoomEvents)
						if err != nil {
							log.Fatal(err)
						}
						for i := 0; i < len(jsonRoomEvents); i++ {
							jsonRoomEvents[i].(map[string]interface{})["waiting_room_id"] = waitingRooms[i].ID
						}
						jsonStructData = append(jsonStructData, jsonRoomEvents...)
					}
					resourceCount = len(jsonStructData)
				case "cloudflare_waiting_room_rules":
					waitingRooms, err := api.ListWaitingRooms(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}
					roomRules := []struct {
						ID            string                 `json:"id"`
						WaitingRoomID string                 `json:"waiting_room_id"`
						Rules         []cfv0.WaitingRoomRule `json:"rules"`
					}{}
					for i := 0; i < len(waitingRooms); i++ {
						rules, err := api.ListWaitingRoomRules(context.Background(), cfv0.ZoneIdentifier(zoneID), cfv0.ListWaitingRoomRuleParams{
							WaitingRoomID: waitingRooms[i].ID,
						})
						if err != nil {
							log.Fatal(err)
						}
						roomRules = append(roomRules, struct {
							ID            string                 `json:"id"`
							WaitingRoomID string                 `json:"waiting_room_id"`
							Rules         []cfv0.WaitingRoomRule `json:"rules"`
						}{
							ID:            waitingRooms[i].ID,
							WaitingRoomID: waitingRooms[i].ID,
							Rules:         rules,
						})
					}
					resourceCount = len(roomRules)
					m, err := json.Marshal(roomRules)
					if err != nil {
						log.Fatal(err)
					}
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_waiting_room_settings":
					waitingRoomSettings, err := api.GetWaitingRoomSettings(context.Background(), cfv0.ZoneIdentifier(zoneID))
					if err != nil {
						log.Fatal(err)
					}
					var jsonPayload []cfv0.WaitingRoomSettings
					jsonPayload = append(jsonPayload, waitingRoomSettings)

					resourceCount = 1
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					jsonStructData[0].(map[string]interface{})["id"] = zoneID
					jsonStructData[0].(map[string]interface{})["search_engine_crawler_bypass"] = waitingRoomSettings.SearchEngineCrawlerBypass
				case "cloudflare_workers_kv_namespace":
					jsonPayload, _, err := api.ListWorkersKVNamespaces(context.Background(), identifier, cfv0.ListWorkersKVNamespacesParams{})
					if err != nil {
						log.Fatal(err)
					}
					resourceCount = len(jsonPayload)
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
					resourceCount = len(jsonPayload.Routes)
					m, _ := json.Marshal(jsonPayload.Routes)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					// remap "script_name" to the "script" value.
					for i := 0; i < resourceCount; i++ {
						jsonStructData[i].(map[string]interface{})["script_name"] = jsonStructData[i].(map[string]interface{})["script"]
					}
				case "cloudflare_zone":
					jsonPayload, err := api.ListZones(context.Background())
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					// - remap "zone" to the "name" value
					// - remap "plan" to "legacy_id" value
					// - drop meta and name_servers
					// - pull in the account_id field
					for i := 0; i < resourceCount; i++ {
						jsonStructData[i].(map[string]interface{})["zone"] = jsonStructData[i].(map[string]interface{})["name"]
						jsonStructData[i].(map[string]interface{})["plan"] = jsonStructData[i].(map[string]interface{})["plan"].(map[string]interface{})["legacy_id"].(string)
						jsonStructData[i].(map[string]interface{})["meta"] = nil
						jsonStructData[i].(map[string]interface{})["name_servers"] = nil
						jsonStructData[i].(map[string]interface{})["status"] = nil
						jsonStructData[i].(map[string]interface{})["account_id"] = jsonStructData[i].(map[string]interface{})["account"].(map[string]interface{})["id"].(string)
					}
				case "cloudflare_zone_lockdown":
					jsonPayload, _, err := api.ListZoneLockdowns(context.Background(), identifier, cfv0.LockdownListParams{})
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

				case "cloudflare_zone_settings_override":
					jsonPayload, err := api.ZoneSettings(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = 1
					m, _ := json.Marshal(jsonPayload.Result)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					zoneSettingsStruct := make(map[string]interface{})
					for _, data := range jsonStructData {
						keyName := data.(map[string]interface{})["id"].(string)
						value := data.(map[string]interface{})["value"]
						zoneSettingsStruct[keyName] = value
					}

					// Remap all settings under "settings" block as well as some of the
					// attributes that are not 1:1 with the API.
					for i := 0; i < resourceCount; i++ {
						jsonStructData[i].(map[string]interface{})["id"] = zoneID
						jsonStructData[i].(map[string]interface{})["settings"] = zoneSettingsStruct

						// zero RTT
						jsonStructData[i].(map[string]interface{})["settings"].(map[string]interface{})["zero_rtt"] = jsonStructData[i].(map[string]interface{})["settings"].(map[string]interface{})["0rtt"]

						// Mobile subdomain redirects
						if jsonStructData[i].(map[string]interface{})["settings"].(map[string]interface{})["mobile_redirect"].(map[string]interface{})["status"] == "off" {
							jsonStructData[i].(map[string]interface{})["settings"].(map[string]interface{})["mobile_redirect"] = nil
						}

						// HSTS
						jsonStructData[i].(map[string]interface{})["settings"].(map[string]interface{})["security_header"].(map[string]interface{})["enabled"] = jsonStructData[i].(map[string]interface{})["settings"].(map[string]interface{})["security_header"].(map[string]interface{})["strict_transport_security"].(map[string]interface{})["enabled"]
						jsonStructData[i].(map[string]interface{})["settings"].(map[string]interface{})["security_header"].(map[string]interface{})["include_subdomains"] = jsonStructData[i].(map[string]interface{})["settings"].(map[string]interface{})["security_header"].(map[string]interface{})["strict_transport_security"].(map[string]interface{})["include_subdomains"]
						jsonStructData[i].(map[string]interface{})["settings"].(map[string]interface{})["security_header"].(map[string]interface{})["max_age"] = jsonStructData[i].(map[string]interface{})["settings"].(map[string]interface{})["security_header"].(map[string]interface{})["strict_transport_security"].(map[string]interface{})["max_age"]
						jsonStructData[i].(map[string]interface{})["settings"].(map[string]interface{})["security_header"].(map[string]interface{})["preload"] = jsonStructData[i].(map[string]interface{})["settings"].(map[string]interface{})["security_header"].(map[string]interface{})["strict_transport_security"].(map[string]interface{})["preload"]
						jsonStructData[i].(map[string]interface{})["settings"].(map[string]interface{})["security_header"].(map[string]interface{})["nosniff"] = jsonStructData[i].(map[string]interface{})["settings"].(map[string]interface{})["security_header"].(map[string]interface{})["strict_transport_security"].(map[string]interface{})["nosniff"]

						// tls_1_2_only is deprecated in favour of min_tls
						jsonStructData[i].(map[string]interface{})["settings"].(map[string]interface{})["tls_1_2_only"] = nil
					}
				case "cloudflare_tiered_cache":
					tieredCache, err := api.GetTieredCache(context.Background(), &cfv0.ResourceContainer{Identifier: zoneID})
					if err != nil {
						log.Fatal(err)
					}
					var jsonPayload []cfv0.TieredCache
					jsonPayload = append(jsonPayload, tieredCache)

					resourceCount = 1
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					jsonStructData[0].(map[string]interface{})["id"] = zoneID
					jsonStructData[0].(map[string]interface{})["cache_type"] = tieredCache.Type.String()
				default:
					fmt.Fprintf(cmd.OutOrStderr(), "%q is not yet supported for automatic generation", resourceType)
					return
				}
			}

			log.Debugf("found %d resources to write out for %q", resourceCount, resourceType)

			// If we don't have any resources to generate, just bail out early.
			if resourceCount == 0 {
				fmt.Fprintf(cmd.OutOrStderr(), "no resources of type %q found to generate", resourceType)
				return
			}

			f := hclwrite.NewEmptyFile()
			rootBody := f.Body()
			for i := 0; i < resourceCount; i++ {
				structData := jsonStructData[i].(map[string]interface{})

				resourceID := ""
				if os.Getenv("USE_STATIC_RESOURCE_IDS") == "true" {
					resourceID = "terraform_managed_resource"
				} else {
					id := ""
					switch structData["id"].(type) {
					case float64:
						id = fmt.Sprintf("%f", structData["id"].(float64))
					default:
						id = structData["id"].(string)
					}

					resourceID = fmt.Sprintf("terraform_managed_resource_%s", id)
				}
				resource := rootBody.AppendNewBlock("resource", []string{resourceType, resourceID}).Body()

				if r == nil {
					log.Fatalf("failed to find %q in the initialized provider schema", resourceType)
				}

				sortedBlockAttributes := make([]string, 0, len(r.Block.Attributes))
				for k := range r.Block.Attributes {
					sortedBlockAttributes = append(sortedBlockAttributes, k)
				}
				sort.Strings(sortedBlockAttributes)

				// Block attributes are for any attributes where assignment is involved.
				for _, attrName := range sortedBlockAttributes {
					// Don't bother outputting the ID for the resource as that is only for
					// internal use (such as importing state).
					if attrName == "id" {
						continue
					}

					// No need to output computed attributes that are also not
					// optional.
					if r.Block.Attributes[attrName].Computed && !r.Block.Attributes[attrName].Optional {
						continue
					}
					if attrName == "account_id" && accountID != "" {
						writeAttrLine(attrName, accountID, "", resource)
						continue
					}

					if attrName == "zone_id" && zoneID != "" && accountID == "" {
						writeAttrLine(attrName, zoneID, "", resource)
						continue
					}

					ty := r.Block.Attributes[attrName].AttributeType
					switch {
					case ty.IsPrimitiveType():
						switch ty {
						case cty.String, cty.Bool, cty.Number:
							writeAttrLine(attrName, structData[attrName], "", resource)
							delete(structData, attrName)
						default:
							log.Debugf("unexpected primitive type %q", ty.FriendlyName())
						}
					case ty.IsCollectionType():
						switch {
						case ty.IsListType(), ty.IsSetType(), ty.IsMapType():
							writeAttrLine(attrName, structData[attrName], "", resource)
							delete(structData, attrName)
						default:
							log.Debugf("unexpected collection type %q", ty.FriendlyName())
						}
					case ty.IsTupleType():
						fmt.Printf("tuple found. attrName %s\n", attrName)
					case ty.IsObjectType():
						fmt.Printf("object found. attrName %s\n", attrName)
					default:
						log.Debugf("attribute %q has not been generated", attrName)
					}
				}

				processBlocks(r.Block, jsonStructData[i].(map[string]interface{}), resource, "")
				f.Body().AppendNewline()
			}

			tfOutput := string(hclwrite.Format(f.Bytes()))
			fmt.Fprint(cmd.OutOrStdout(), tfOutput)
		}
	}
}
