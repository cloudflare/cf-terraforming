package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"sort"
	"strings"

	cfv0 "github.com/cloudflare/cloudflare-go"
	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"github.com/zclconf/go-cty/cty"
)

var (
	resourceType    string
	resourceIDFlags []string

	generateCmd = &cobra.Command{
		Use:    "generate",
		Short:  "Fetch resources from the Cloudflare API and generate the respective Terraform stanzas",
		Run:    generateResources(),
		PreRun: sharedPreRun,
	}

	deprecatedResources = []string{"cloudflare_firewall_rule"}
)

func init() {
	rootCmd.AddCommand(generateCmd)
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
		log.WithFields(logrus.Fields{
			"directory": workingDir,
		}).Debug("initializing Terraform")
		tf, err := tfexec.NewTerraform(workingDir, execPath)
		if err != nil {
			log.Fatal(err)
		}

		_, providerVersion, err := tf.Version(context.Background(), true)
		if err != nil {
			log.Fatalf("failed to retrieve terraform and provider version information: %s", err)
		}

		var registryPath string
		for provider := range providerVersion {
			if strings.Contains(provider, "/cloudflare/cloudflare") {
				registryPath = provider
				continue
			}
		}

		detectedVersion, ok := providerVersion[registryPath]
		if !ok {
			log.WithFields(logrus.Fields{
				"available_registries": providerVersion,
			}).Fatal("failed to find registry")
		}

		providerVersionString := detectedVersion.String()
		log.WithFields(logrus.Fields{
			"version":  providerVersionString,
			"registry": registryPath,
		}).Debug("detected provider")

		log.Debug("reading Terraform schema")
		ps, err := tf.ProvidersSchema(context.Background())
		if err != nil {
			log.Fatal("failed to read provider schema", err)
		}

		s := ps.Schemas[registryPath]
		if s == nil {
			log.Fatal("failed to detect provider installation")
		}

		resources := strings.Split(resourceType, ",")
		for _, resourceType := range resources {
			r := s.ResourceSchemas[resourceType]
			log.WithFields(logrus.Fields{
				"resource": resourceType,
			}).Debug("reading and building resource")
			if (r != nil && r.Block != nil && r.Block.Deprecated) || slices.Contains(deprecatedResources, resourceType) {
				log.Warnf(fmt.Sprintf("resource %s is deprecated. The terraform config might not be generated.", resourceType))
			}

			// Initialise `resourceCount` outside of the switch for supported resources
			// to allow it to be referenced further down in the loop that outputs the
			// newly generated resources.
			resourceCount := 0
			var jsonStructData []interface{}

			// The ruleset API has many gotchas that are accounted for in how we build
			// the 'response' object that feeds into the HCL generation, and it's difficult
			// to ensure the same compatability using the generated SDK.
			useOldSDK := resourceType == "cloudflare_ruleset"

			if strings.HasPrefix(providerVersionString, "5") && !useOldSDK {
				resourceIDsMap := make(map[string][]string)
				if isSupportedPathParam(resources, resourceType) {
					resourceIDsMap = getResourceMappings()

					ids, ok := resourceIDsMap[resourceType]
					if ok && len(ids) == 0 {
						log.Fatalf("No resource IDs defined in Terraform for resource %s", resourceType)
					}
				}

				if resourceToEndpoint[resourceType]["list"] == "" && resourceToEndpoint[resourceType]["get"] == "" {
					log.WithFields(logrus.Fields{
						"resource": resourceType,
					}).Debug("did not find API endpoint. does it exist in the mapping?")
					continue
				}

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
				if strings.Contains(endpoint, "{accounts_or_zones}") {
					if accountID != "" {
						endpoint = strings.Replace(endpoint, "/{accounts_or_zones}/{account_or_zone_id}/", "/accounts/{account_id}/", 1)
					} else {
						endpoint = strings.Replace(endpoint, "/{accounts_or_zones}/{account_or_zone_id}/", "/zones/{zone_id}/", 1)
					}
				}

				// replace the URL placeholders with the actual values we have.
				placeholderReplacer := strings.NewReplacer("{account_id}", accountID, "{zone_id}", zoneID)
				endpoint = placeholderReplacer.Replace(endpoint)

				pathParams, ok := resourceIDsMap[resourceType]
				if ok && len(pathParams) > 0 {
					endpoints := replacePathParams(pathParams, endpoint, resourceType)
					jsonStructData, err = GetAPIResponse(result, pathParams, endpoints...)
					if err != nil {
						log.Infof("error getting API response for resource %s: %s", resourceType, err)
						continue
					}
					resourceCount = len(jsonStructData)
				} else {
					jsonStructData, err = GetAPIResponse(result, pathParams, endpoint)
					if err != nil {
						log.Infof("error getting API response for resource %s: %s", resourceType, err)
						continue
					}
					resourceCount = len(jsonStructData)
				}
			} else {
				var identifier *cfv0.ResourceContainer
				if accountID != "" {
					identifier = cfv0.AccountIdentifier(accountID)
				} else {
					identifier = cfv0.ZoneIdentifier(zoneID)
				}

				switch resourceType {
				case "cloudflare_access_application":
					jsonPayload, _, err := apiV0.ListAccessApplications(context.Background(), identifier, cfv0.ListAccessApplicationsParams{})
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
					jsonPayload, _, err := apiV0.ListAccessGroups(context.Background(), identifier, cfv0.ListAccessGroupsParams{})
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
					jsonPayload, _, err := apiV0.ListAccessIdentityProviders(context.Background(), identifier, cfv0.ListAccessIdentityProvidersParams{})
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
					jsonPayload, _, err := apiV0.ListAccessServiceTokens(context.Background(), identifier, cfv0.ListAccessServiceTokensParams{})
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
					jsonPayload, _, err := apiV0.ListAccessMutualTLSCertificates(context.Background(), identifier, cfv0.ListAccessMutualTLSCertificatesParams{})
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
						jsonPayload, err := apiV0.ListAccountAccessRules(context.Background(), accountID, cfv0.AccessRule{}, 1)
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
						jsonPayload, err := apiV0.ListZoneAccessRules(context.Background(), zoneID, cfv0.AccessRule{}, 1)
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
					jsonPayload, _, err := apiV0.AccountMembers(context.Background(), accountID, cfv0.PaginationOptions{})
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

					argoSmartRouting, err := apiV0.ArgoSmartRouting(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}
					jsonPayload = append(jsonPayload, argoSmartRouting)

					argoTieredCaching, err := apiV0.ArgoTieredCaching(context.Background(), zoneID)
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
					apiShieldConfig, _, err := apiV0.GetAPIShieldConfiguration(context.Background(), identifier)
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
						res, err := apiV0.ListUserAgentRules(context.Background(), zoneID, page)
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
					botManagement, err := apiV0.GetBotManagement(context.Background(), identifier)
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
					jsonPayload, err := apiV0.ListPrefixes(context.Background(), accountID)
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
					jsonPayload, err := apiV0.ListCertificatePacks(context.Background(), zoneID)
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
						jsonPayload, err := apiV0.CustomPages(context.Background(), &acc)
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
						jsonPayload, err := apiV0.CustomPages(context.Background(), &zo)
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
					apiCall, err := apiV0.CustomHostnameFallbackOrigin(context.Background(), zoneID)
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
					jsonPayload, _, err := apiV0.Filters(context.Background(), identifier, cfv0.FilterListParams{})
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
					jsonPayload, _, err := apiV0.FirewallRules(context.Background(), identifier, cfv0.FirewallRuleListParams{})
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
					jsonPayload, _, err := apiV0.CustomHostnames(context.Background(), zoneID, 1, cfv0.CustomHostname{})
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
					jsonPayload, err := apiV0.ListSSL(context.Background(), zoneID)
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
					jsonPayload, err := apiV0.Healthchecks(context.Background(), zoneID)
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
					jsonPayload, err := apiV0.ListLists(context.Background(), identifier, cfv0.ListListsParams{})
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

						listItems, err := apiV0.ListListItems(context.Background(), identifier, cfv0.ListListItemsParams{ID: listID})
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
					jsonPayload, err := apiV0.ListLoadBalancers(context.Background(), identifier, cfv0.ListLoadBalancerParams{})
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
					jsonPayload, err := apiV0.ListLoadBalancerPools(context.Background(), identifier, cfv0.ListLoadBalancerPoolParams{})
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
					jsonPayload, err := apiV0.ListLoadBalancerMonitors(context.Background(), identifier, cfv0.ListLoadBalancerMonitorParams{})
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
					jsonPayload, err := apiV0.ListLogpushJobs(context.Background(), identifier, cfv0.ListLogpushJobsParams{})
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
					jsonPayload, err := apiV0.ListZoneManagedHeaders(context.Background(), cfv0.ResourceIdentifier(zoneID), cfv0.ListManagedHeadersParams{Status: "enabled"})
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
					jsonPayload, err := apiV0.ListOriginCACertificates(context.Background(), cfv0.ListOriginCertificatesParams{ZoneID: zoneID})
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
					jsonPayload, err := apiV0.ListPageRules(context.Background(), zoneID)
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
					jsonPayload, err := apiV0.ListAllRateLimits(context.Background(), zoneID)
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
					jsonPayload, _, err := apiV0.ListDNSRecords(context.Background(), identifier, cfv0.ListDNSRecordsParams{})
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					zone, _ := apiV0.ZoneDetails(context.Background(), identifier.Identifier)

					for i := 0; i < resourceCount; i++ {
						// Drop the proxiable values as they are not usable
						jsonStructData[i].(map[string]interface{})["proxiable"] = nil
						jsonStructData[i].(map[string]interface{})["value"] = nil

						if jsonStructData[i].(map[string]interface{})["name"].(string) != zone.Name {
							jsonStructData[i].(map[string]interface{})["name"] = strings.ReplaceAll(jsonStructData[i].(map[string]interface{})["name"].(string), "."+zone.Name, "")
						}
					}
				case "cloudflare_ruleset":
					jsonPayload, err := apiV0.ListRulesets(context.Background(), identifier, cfv0.ListRulesetsParams{})
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
						ruleset, _ := apiV0.GetRuleset(context.Background(), identifier, rule.ID)
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

					if strings.HasPrefix(providerVersionString, "5") {
						for i := 0; i < resourceCount; i++ {
							rules := jsonStructData[i].(map[string]interface{})["rules"]
							if rules != nil {
								for ruleCounter := range rules.([]interface{}) {
									rules.([]interface{})[ruleCounter].(map[string]interface{})["id"] = nil
								}
							}
						}
						continue
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
					jsonPayload, err := apiV0.SpectrumApplications(context.Background(), zoneID)
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
					jsonPayload, _, err := apiV0.ListTeamsLists(context.Background(), identifier, cfv0.ListTeamListsParams{})
					if err != nil {
						log.Fatal(err)
					}
					// get items for the lists and add it the specific list struct
					for i, TeamsList := range jsonPayload {
						items_struct, _, err := apiV0.ListTeamsListItems(
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
					jsonPayload, _, err := apiV0.TeamsLocations(context.Background(), accountID)
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
					jsonPayload, _, err := apiV0.TeamsProxyEndpoints(context.Background(), accountID)
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
					jsonPayload, err := apiV0.TeamsRules(context.Background(), accountID)
					if err != nil {
						log.Fatal(err)
					}
					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}

					// flatten add_headers of rule setting to a string
					for i := 0; i < resourceCount; i++ {
						ruleSettings, ok := jsonStructData[i].(map[string]interface{})["rule_settings"].(map[string]interface{})
						if ok {
							addHeaders, ok := ruleSettings["add_headers"].(map[string]interface{})
							if ok {
								for k, v := range addHeaders {
									headerValues := v.([]interface{})
									headerString := ""
									for _, headerValue := range headerValues {
										headerString += strings.Join([]string{headerValue.(string)}, ",")
									}
									addHeaders[k] = headerString
								}
							}
						}
						// check for empty descriptions
						if jsonStructData[i].(map[string]interface{})["description"] == "" {
							jsonStructData[i].(map[string]interface{})["description"] = "default"
						}
					}
				case "cloudflare_tunnel":
					log.Debug("only requesting the first 1000 active Cloudflare Tunnels due to the service not providing correct pagination responses")
					jsonPayload, _, err := apiV0.ListTunnels(
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
						secret, err := apiV0.GetTunnelToken(
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
					jsonPayload, _, err := apiV0.ListTurnstileWidgets(context.Background(), identifier, cfv0.ListTurnstileWidgetParams{})
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
					jsonPayload, err := apiV0.URLNormalizationSettings(context.Background(), &cfv0.ResourceContainer{Identifier: zoneID, Level: cfv0.ZoneRouteLevel})
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
					jsonPayload, err := apiV0.ListWaitingRooms(context.Background(), zoneID)
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
					waitingRooms, err := apiV0.ListWaitingRooms(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}
					for i := 0; i < len(waitingRooms); i++ {
						roomEvents, err := apiV0.ListWaitingRoomEvents(context.Background(), zoneID, waitingRooms[i].ID)
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
					waitingRooms, err := apiV0.ListWaitingRooms(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}
					roomRules := []struct {
						ID            string                 `json:"id"`
						WaitingRoomID string                 `json:"waiting_room_id"`
						Rules         []cfv0.WaitingRoomRule `json:"rules"`
					}{}
					for i := 0; i < len(waitingRooms); i++ {
						rules, err := apiV0.ListWaitingRoomRules(context.Background(), cfv0.ZoneIdentifier(zoneID), cfv0.ListWaitingRoomRuleParams{
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
					waitingRoomSettings, err := apiV0.GetWaitingRoomSettings(context.Background(), cfv0.ZoneIdentifier(zoneID))
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
					jsonPayload, _, err := apiV0.ListWorkersKVNamespaces(context.Background(), identifier, cfv0.ListWorkersKVNamespacesParams{})
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
					jsonPayload, err := apiV0.ListWorkerRoutes(context.Background(), identifier, cfv0.ListWorkerRoutesParams{})
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
					jsonPayload, err := apiV0.ListZones(context.Background())
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
					jsonPayload, _, err := apiV0.ListZoneLockdowns(context.Background(), identifier, cfv0.LockdownListParams{})
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
					jsonPayload, err := apiV0.ZoneSettings(context.Background(), zoneID)
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
					tieredCache, err := apiV0.GetTieredCache(context.Background(), &cfv0.ResourceContainer{Identifier: zoneID})
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

			log.WithFields(logrus.Fields{
				"count":    resourceCount,
				"resource": resourceType,
			}).Debug("generating resource output")

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
					if resourceCount == 1 {
						resourceID = "terraform_managed_resource"
					} else {
						resourceID = fmt.Sprintf("terraform_managed_resource_%d", i)
					}
				} else {
					id := ""
					switch structData["id"].(type) {
					case float64:
						id = fmt.Sprintf("%f", structData["id"].(float64))
					default:
						if structData["id"] == nil {
							if accountID != "" {
								id = accountID
							}

							if zoneID != "" {
								id = zoneID
							}
						} else {
							id = structData["id"].(string)
						}
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

			postProcess(f, resourceType)
			tfOutput := string(hclwrite.Format(f.Bytes()))
			fmt.Fprint(cmd.OutOrStdout(), tfOutput)
		}
	}
}

func processCustomCasesV5(response *[]interface{}, resourceType string, pathParam string) {
	resourceCount := len(*response)
	switch resourceType {
	case "cloudflare_managed_transforms":
		// remap email and role_ids into the right structure and remove policies
		for i := 0; i < resourceCount; i++ {
			for j := range (*response)[i].(map[string]interface{})["managed_request_headers"].([]interface{}) {
				delete((*response)[i].(map[string]interface{})["managed_request_headers"].([]interface{})[j].(map[string]interface{}), "has_conflict")
			}
			for j := range (*response)[i].(map[string]interface{})["managed_response_headers"].([]interface{}) {
				delete((*response)[i].(map[string]interface{})["managed_response_headers"].([]interface{})[j].(map[string]interface{}), "has_conflict")
			}
		}
	case "cloudflare_r2_bucket":
		finalResponse := make([]interface{}, 0)
		r := *response
		for i := 0; i < resourceCount; i++ {
			buckets := r[i].(map[string]interface{})["buckets"]
			bucketObjects := make([]interface{}, len(buckets.([]interface{})))
			for j := range buckets.([]interface{}) {
				b := buckets.([]interface{})[j]
				bucketObjects[j] = b
			}
			finalResponse = append(finalResponse, bucketObjects...)
		}
		*response = make([]interface{}, len(finalResponse))
		for i := range finalResponse {
			(*response)[i] = finalResponse[i]
		}
	case "cloudflare_account_member":
		// remap email and role_ids into the right structure and remove policies
		for i := 0; i < resourceCount; i++ {
			delete((*response)[i].(map[string]interface{}), "policies")
			(*response)[i].(map[string]interface{})["email"] = (*response)[i].(map[string]interface{})["user"].(map[string]interface{})["email"]
			roleIDs := []string{}
			for _, role := range (*response)[i].(map[string]interface{})["roles"].([]interface{}) {
				roleIDs = append(roleIDs, role.(map[string]interface{})["id"].(string))
			}
			(*response)[i].(map[string]interface{})["roles"] = roleIDs
		}
	case "cloudflare_content_scanning_expression":
		// wrap the response in 'body' for tf
		for i := 0; i < resourceCount; i++ {
			payload := (*response)[i].(map[string]interface{})["payload"]
			(*response)[i].(map[string]interface{})["body"] = []interface{}{map[string]interface{}{
				"payload": payload,
			}}
		}
	case "cloudflare_zero_trust_device_default_profile_local_domain_fallback":
		// wrap the response in 'domains' for tf
		for i := 0; i < resourceCount; i++ {
			do := make(map[string]interface{})
			do["domains"] = []interface{}{(*response)[i]}
			(*response)[i] = do
		}
	case "cloudflare_zero_trust_dex_test":
		// remove the nesting under 'dex_test'
		finalResponse := make([]interface{}, 0)
		r := *response
		for i := 0; i < resourceCount; i++ {
			dexTests := r[i].(map[string]interface{})["dex_tests"]
			dtObjects := make([]interface{}, len(dexTests.([]interface{})))
			for j := range dexTests.([]interface{}) {
				dt := dexTests.([]interface{})[j]
				dtObjects[j] = dt
			}
			finalResponse = append(finalResponse, dtObjects...)
		}
		*response = make([]interface{}, len(finalResponse))
		for i := range finalResponse {
			(*response)[i] = finalResponse[i]
		}
	case "cloudflare_zero_trust_gateway_settings":
		for i := 0; i < resourceCount; i++ {
			settings, ok := (*response)[i].(map[string]interface{})["settings"]
			if !ok {
				return
			}
			customCert, ok := settings.(map[string]interface{})["custom_certificate"]
			if ok {
				delete(customCert.(map[string]interface{}), "binding_status")
				delete(customCert.(map[string]interface{}), "expires_on")
				delete(customCert.(map[string]interface{}), "updated_at")
			}
		}
	case "cloudflare_page_rule":
		for i := 0; i < resourceCount; i++ {
			(*response)[i].(map[string]interface{})["target"] = (*response)[i].(map[string]interface{})["targets"].([]interface{})[0].(map[string]interface{})["constraint"].(map[string]interface{})["value"]
			(*response)[i].(map[string]interface{})["actions"] = flattenAttrMap((*response)[i].(map[string]interface{})["actions"].([]interface{}))

			// Have to remap the cache_ttl_by_status to conform to Terraform's more human-friendly structure.
			if cache, ok := (*response)[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_ttl_by_status"].(map[string]interface{}); ok {
				cacheTtlByStatus := []map[string]interface{}{}

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

					cacheTtlByStatus = append(cacheTtlByStatus, elem)
				}

				sort.SliceStable(cacheTtlByStatus, func(i int, j int) bool {
					return cacheTtlByStatus[i]["codes"].(string) < cacheTtlByStatus[j]["codes"].(string)
				})

				(*response)[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_ttl_by_status"] = cacheTtlByStatus
			}

			// Remap cache_key_fields.query_string.include & .exclude wildcards (not in an array) to the appropriate "ignore" field value in Terraform.
			if c, ok := (*response)[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_key_fields"].(map[string]interface{}); ok {
				if s, sok := c["query_string"].(map[string]interface{})["include"].(string); sok && s == "*" {
					(*response)[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_key_fields"].(map[string]interface{})["query_string"].(map[string]interface{})["include"] = nil
					(*response)[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_key_fields"].(map[string]interface{})["query_string"].(map[string]interface{})["ignore"] = false
				}
				if s, sok := c["query_string"].(map[string]interface{})["exclude"].(string); sok && s == "*" {
					(*response)[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_key_fields"].(map[string]interface{})["query_string"].(map[string]interface{})["exclude"] = nil
					(*response)[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_key_fields"].(map[string]interface{})["query_string"].(map[string]interface{})["ignore"] = true
				}
			}
		}
	case "cloudflare_zero_trust_access_short_lived_certificate":
		// map id under app_id
		for i := 0; i < resourceCount; i++ {
			appID := (*response)[i].(map[string]interface{})["id"]
			(*response)[i].(map[string]interface{})["app_id"] = appID
		}
	case "cloudflare_zone_setting":
		for i := 0; i < resourceCount; i++ {
			(*response)[i].(map[string]interface{})["setting_id"] = (*response)[i].(map[string]interface{})["id"]
		}
	case "cloudflare_hostname_tls_setting":
		for i := 0; i < resourceCount; i++ {
			(*response)[i].(map[string]interface{})["setting_id"] = pathParam
		}
	case "cloudflare_registrar_domain":
		for i := 0; i < resourceCount; i++ {
			(*response)[i].(map[string]interface{})["domain_name"] = (*response)[i].(map[string]interface{})["name"]
		}
	case "cloudflare_r2_managed_domain":
		for i := 0; i < resourceCount; i++ {
			(*response)[i].(map[string]interface{})["bucket_name"] = pathParam
		}
	case "cloudflare_r2_custom_domain":
		finalResponse := make([]interface{}, 0)
		r := *response
		for i := 0; i < resourceCount; i++ {
			domains := r[i].(map[string]interface{})["domains"]
			bucketObjects := make([]interface{}, len(domains.([]interface{})))
			for j := range domains.([]interface{}) {
				b := domains.([]interface{})[j]
				b.(map[string]interface{})["bucket_name"] = pathParam
				b.(map[string]interface{})["zone_id"] = b.(map[string]interface{})["zoneId"]
				bucketObjects[j] = b
			}
			finalResponse = append(finalResponse, bucketObjects...)
		}
		*response = make([]interface{}, len(finalResponse))
		for i := range finalResponse {
			(*response)[i] = finalResponse[i]
		}
	case "cloudflare_pages_domain":
		for i := 0; i < resourceCount; i++ {
			(*response)[i].(map[string]interface{})["project_name"] = pathParam
		}
	case "cloudflare_list_item":
		for i := 0; i < resourceCount; i++ {
			(*response)[i].(map[string]interface{})["list_id"] = (*response)[i].(map[string]interface{})["id"]
		}
	case "cloudflare_api_shield_schema":
		for i := 0; i < resourceCount; i++ {
			(*response)[i].(map[string]interface{})["file"] = (*response)[i].(map[string]interface{})["source"]
		}
	case "cloudflare_api_shield_discovery_operation":
		for i := 0; i < resourceCount; i++ {
			(*response)[i].(map[string]interface{})["operation_id"] = (*response)[i].(map[string]interface{})["id"]
		}
	case "cloudflare_zero_trust_dlp_predefined_profile":
		for i := 0; i < resourceCount; i++ {
			(*response)[i].(map[string]interface{})["profile_id"] = pathParam
		}
	case "cloudflare_zero_trust_access_identity_provider":
		for i := 0; i < resourceCount; i++ {
			cfg, ok := (*response)[i].(map[string]interface{})["config"]
			if ok {
				delete(cfg.(map[string]interface{}), "redirect_url")
			}
			scimCFG, ok := (*response)[i].(map[string]interface{})["scim_config"]
			if ok {
				delete(scimCFG.(map[string]interface{}), "scim_base_url")
			}
		}
	case "cloudflare_zero_trust_access_custom_page":
		// fetch each object one by one to get 'custom_html' field.
		endpointFMT := resourceToEndpoint[resourceType]["get"]
		placeholderReplacer := strings.NewReplacer("{account_id}", accountID)
		endpointFMT = placeholderReplacer.Replace(endpointFMT)
		for i := 0; i < resourceCount; i++ {
			uid, ok := (*response)[i].(map[string]interface{})["uid"]
			if !ok {
				continue
			}
			endpoint := strings.Replace(endpointFMT, "{custom_page_id}", uid.(string), 1)
			result := new(http.Response)
			err := api.Get(context.Background(), endpoint, nil, &result)
			if err != nil {
				var apierr *cloudflare.Error
				if errors.As(err, &apierr) {
					if apierr.StatusCode == http.StatusNotFound {
						log.WithFields(logrus.Fields{
							"resource": resourceType,
							"endpoint": endpoint,
						}).Debug("no resources found")
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
			customHTML := gjson.Get(value.Raw, "custom_html")
			if value.Type == gjson.Null {
				continue
			}
			(*response)[i].(map[string]interface{})["custom_html"] = customHTML.String()
		}
	case "cloudflare_web_analytics_rule":
		finalResponse := make([]interface{}, 0)
		r := *response
		for i := 0; i < resourceCount; i++ {
			rules := r[i].(map[string]interface{})["rules"]
			ruleObjects := make([]interface{}, len(rules.([]interface{})))
			for j := range rules.([]interface{}) {
				b := rules.([]interface{})[j]
				b.(map[string]interface{})["ruleset_id"] = pathParam
				ruleObjects[j] = b
			}
			finalResponse = append(finalResponse, ruleObjects...)
		}
		*response = make([]interface{}, len(finalResponse))
		for i := range finalResponse {
			(*response)[i] = finalResponse[i]
		}
	}
}

func unMarshallJSONStructData(modifiedJSONString string) ([]interface{}, error) {
	var data interface{}
	err := json.Unmarshal([]byte(modifiedJSONString), &data)
	if err != nil {
		return nil, err
	}
	if dataSlice, ok := data.([]interface{}); ok {
		return dataSlice, nil
	}
	return []interface{}{data}, nil
}

// postProcess allows you to perform additional actions on the generated hcl.
func postProcess(f *hclwrite.File, resourceType string) {
	switch resourceType {
	case "cloudflare_stream_live_input", "cloudflare_stream":
		addJSONEncode(f, "meta")
	}
}

// addJSONEncode wraps a hcl block with the jsonencode function.
func addJSONEncode(f *hclwrite.File, attributeName string) {
	for _, block := range f.Body().Blocks() {
		if block.Type() != "resource" {
			continue
		}
		if len(block.Labels()) < 1 {
			continue
		}
		if block.Labels()[0] != resourceType {
			continue
		}
		body := block.Body()
		attr := body.GetAttribute(attributeName)
		if attr == nil {
			continue
		}
		exprTokens := attr.Expr().BuildTokens(nil)
		exprText := string(exprTokens.Bytes())

		trimmed := strings.TrimSpace(exprText)
		// Wrap the attribute with jsonencode
		if len(trimmed) > 0 && trimmed[0] == '{' {
			body.RemoveAttribute(attributeName)
			newTokens := hclwrite.Tokens{}
			fnStart := &hclwrite.Token{
				Type:  hclsyntax.TokenIdent,
				Bytes: []byte("jsonencode("),
			}
			newTokens = append(newTokens, fnStart)
			newTokens = append(newTokens, exprTokens...)
			fnEnd := &hclwrite.Token{
				Type:  hclsyntax.TokenCParen,
				Bytes: []byte(")"),
			}
			newTokens = append(newTokens, fnEnd)
			body.SetAttributeRaw(attributeName, newTokens)
		}
	}
}

func GetAPIResponse(result *http.Response, pathParams []string, endpoints ...string) ([]interface{}, error) {
	var jsonStructData, results []interface{}
	for i, endpoint := range endpoints {
		err := api.Get(context.Background(), endpoint, nil, &result)
		if err != nil {
			var apierr *cloudflare.Error
			if errors.As(err, &apierr) {
				if apierr.StatusCode == http.StatusNotFound {
					log.WithFields(logrus.Fields{
						"resource": resourceType,
						"endpoint": endpoint,
					}).Debug("no resources found")
					return nil, err
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
			return nil, errors.New("no result found")
		}

		modifiedJSON := modifyResponsePayload(resourceType, value)
		jsonStructData, err = unMarshallJSONStructData(modifiedJSON)
		if err != nil {
			log.Fatalf("failed to unmarshal result: %s", err)
		}

		param := ""
		if len(pathParams) > 0 {
			param = pathParams[i]
		}
		processCustomCasesV5(&jsonStructData, resourceType, param)
		results = append(results, jsonStructData...)
	}
	return results, nil
}

func isSupportedPathParam(resources []string, rType string) bool {
	_, ok := settingsMap[rType]
	if !ok {
		return false
	}
	return slices.Contains(resources, rType)
}

func replacePathParams(params []string, endpoint string, rType string) []string {
	endpoints := make([]string, 0)
	var placeholder string
	switch rType {
	case "cloudflare_zone_setting", "cloudflare_hostname_tls_setting":
		placeholder = "{setting_id}"
	case "cloudflare_waiting_room_event":
		placeholder = "{waiting_room_id}"
	case "cloudflare_r2_managed_domain", "cloudflare_r2_custom_domain":
		placeholder = "{bucket_name}"
	case "cloudflare_pages_domain":
		placeholder = "{project_name}"
	case "cloudflare_list_item":
		placeholder = "{list_id}"
	case "cloudflare_zero_trust_dlp_predefined_profile":
		placeholder = "{profile_id}"
	case "cloudflare_web_analytics_rule":
		placeholder = "{ruleset_id}"

	default:
		return endpoints
	}
	for _, id := range params {
		endpoints = append(endpoints, strings.Clone(strings.NewReplacer(placeholder, id).Replace(endpoint)))
	}
	return endpoints
}
