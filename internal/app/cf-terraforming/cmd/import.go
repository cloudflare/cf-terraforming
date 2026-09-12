package cmd

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	cfv0 "github.com/cloudflare/cloudflare-go"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zclconf/go-cty/cty"
)

// resourceImportStringFormats contains a mapping of the resource type to the
// composite ID that is compatible with performing an import.
// resourceIDFieldOverrides names the field holding a resource's own identifier
// for API responses that do not use "id".
//
// Without these, the entries fall through to the account or zone ID (see the
// fallback in the import loop), which gives every object of the type the same
// import id and the same generated resource name. The result passes a parity
// check -- one import block per resource -- while pointing all of them at one
// object.
var resourceIDFieldOverrides = map[string]string{
	"cloudflare_r2_bucket":                        "name",
	"cloudflare_web_analytics_site":               "site_tag",
	"cloudflare_zero_trust_device_custom_profile": "policy_id",
	"cloudflare_zero_trust_dex_test":              "test_id",
}

// identifierFromResponse returns the resource's own identifier from an API
// response object whose identifier is not named "id", or "" when there is none.
// Shared by generate and import so a resource's name and its import id are
// always derived the same way.
func identifierFromResponse(resourceType string, data map[string]interface{}) string {
	field, ok := resourceIDFieldOverrides[resourceType]
	if !ok {
		return ""
	}
	v, present := data[field]
	if !present || v == nil {
		return ""
	}
	str, isStr := v.(string)
	if !isStr {
		return ""
	}
	return str
}

// scopeAddressedIdentifier returns the scope id a resource is addressed by when
// its endpoint carries only a scope placeholder, or "" when it is not such a
// resource.
//
// Used for naming as well as for the import id: a scope-addressed resource takes
// its object's own identifier from the API (a setting name, say), and naming
// every zone's copy after that produces the same Terraform address in all of
// them. Terraform rejects a repeated address, so the output does not parse even
// before the import ids are considered.
func scopeAddressedIdentifier(endpoint string) string {
	e := endpoint
	if strings.Contains(e, "{accounts_or_zones}") {
		if accountID != "" {
			e = strings.Replace(e, "/{accounts_or_zones}/{account_or_zone_id}/", "/accounts/{account_id}/", 1)
		} else {
			e = strings.Replace(e, "/{accounts_or_zones}/{account_or_zone_id}/", "/zones/{zone_id}/", 1)
		}
	}
	matches := regexp.MustCompile("({[a-z0-9_]*})").FindAllString(e, -1)
	if len(matches) != 1 || !isScopePlaceholder(matches[0]) {
		return ""
	}
	if matches[0] == "{account_id}" {
		return accountID
	}
	return zoneID
}

// isScopePlaceholder reports whether an endpoint placeholder names the account or
// zone the request is scoped to, rather than an object within it.
func isScopePlaceholder(p string) bool {
	return p == "{account_id}" || p == "{zone_id}"
}

// hclSafeIdentifier makes a value safe to embed in a Terraform resource name.
// Identifiers taken from fields such as an R2 bucket name are not restricted to
// the characters HCL allows in an identifier -- a bucket may contain a dot, for
// example -- and a resource name is not a quoted string, so an unsanitised value
// produces a file that will not parse. Previously this position only ever held a
// hex or UUID id, so the problem was unreachable.
func hclSafeIdentifier(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_', r == '-':
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	return b.String()
}

// idFieldHint describes which alternative field was tried, for the error message.
func idFieldHint(resourceType string) string {
	if field, ok := resourceIDFieldOverrides[resourceType]; ok {
		return fmt.Sprintf(" and %q", field)
	}
	return ""
}

var resourceImportStringFormats = map[string]string{
	"cloudflare_access_application":                            ":account_id/:id",
	"cloudflare_access_group":                                  ":account_id/:id",
	"cloudflare_access_rule":                                   ":identifier_type/:identifier_value/:id",
	"cloudflare_account":                                       ":account_id",
	"cloudflare_account_member":                                ":account_id/:id",
	"cloudflare_api_shield_operation":                          ":zone_id/:id",
	"cloudflare_argo":                                          ":zone_id/argo",
	"cloudflare_bot_management":                                ":zone_id",
	"cloudflare_byo_ip_prefix":                                 ":id",
	"cloudflare_certificate_pack":                              ":zone_id/:id",
	"cloudflare_custom_hostname":                               ":zone_id/:id",
	"cloudflare_custom_pages":                                  ":identifier_type/:identifier_value/:id",
	"cloudflare_custom_ssl":                                    ":zone_id/:id",
	"cloudflare_d1_database":                                   ":account_id/:id",
	"cloudflare_dns_firewall":                                  ":account_id/:id",
	"cloudflare_dns_record":                                    ":zone_id/:id",
	"cloudflare_dns_zone_transfers_acl":                        ":account_id/:id",
	"cloudflare_dns_zone_transfers_incoming":                   ":zone_id",
	"cloudflare_dns_zone_transfers_outgoing":                   ":zone_id",
	"cloudflare_dns_zone_transfers_peer":                       ":account_id/:id",
	"cloudflare_dns_zone_transfers_tsig":                       ":account_id/:id",
	"cloudflare_email_routing_address":                         ":account_id/:id",
	"cloudflare_email_routing_catch_all":                       ":zone_id",
	"cloudflare_email_routing_dns":                             ":zone_id",
	"cloudflare_email_routing_rule":                            ":zone_id/:id",
	"cloudflare_email_routing_settings":                        ":zone_id",
	"cloudflare_email_security_block_sender":                   ":account_id/:id",
	"cloudflare_email_security_impersonation_registry":         ":account_id/:id",
	"cloudflare_email_security_trusted_domains":                ":account_id/:id",
	"cloudflare_filter":                                        ":zone_id/:id",
	"cloudflare_firewall_rule":                                 ":zone_id/:id",
	"cloudflare_healthcheck":                                   ":zone_id/:id",
	"cloudflare_ip_list":                                       ":account_id/:id",
	"cloudflare_keyless_certificate":                           ":zone_id/:id",
	"cloudflare_list":                                          ":account_id/:id",
	"cloudflare_load_balancer":                                 ":zone_id/:id",
	"cloudflare_load_balancer_monitor":                         ":account_id/:id",
	"cloudflare_load_balancer_pool":                            ":account_id/:id",
	"cloudflare_managed_transforms":                            ":zone_id",
	"cloudflare_mtls_certificate":                              ":account_id/:id",
	"cloudflare_notification_policy":                           ":account_id/:id",
	"cloudflare_notification_policy_webhooks":                  ":account_id/:id",
	"cloudflare_origin_ca_certificate":                         ":id",
	"cloudflare_page_rule":                                     ":zone_id/:id",
	"cloudflare_page_shield_policy":                            ":zone_id/:id",
	"cloudflare_pages_project":                                 ":account_id/:id",
	"cloudflare_queue":                                         ":account_id/:id",
	"cloudflare_r2_bucket":                                     ":account_id/:id",
	"cloudflare_rate_limit":                                    ":zone_id/:id",
	"cloudflare_record":                                        ":zone_id/:id",
	"cloudflare_regional_hostname":                             ":zone_id/:id",
	"cloudflare_regional_tiered_cache":                         ":zone_id",
	"cloudflare_ruleset":                                       ":identifier_type/:identifier_value/:id",
	"cloudflare_spectrum_application":                          ":zone_id/:id",
	"cloudflare_stream_key":                                    ":account_id",
	"cloudflare_tiered_cache":                                  ":zone_id",
	"cloudflare_total_tls":                                     ":zone_id",
	"cloudflare_tunnel":                                        ":account_id/:id",
	"cloudflare_turnstile_widget":                              ":account_id/:id",
	"cloudflare_url_normalization_settings":                    ":zone_id",
	"cloudflare_waf_override":                                  ":zone_id/:id",
	"cloudflare_waiting_room":                                  ":zone_id/:id",
	"cloudflare_waiting_room_settings":                         ":zone_id",
	"cloudflare_web3_hostname":                                 ":zone_id/:id",
	"cloudflare_web_analytics_site":                            ":account_id/:id",
	"cloudflare_worker_route":                                  ":zone_id/:id",
	"cloudflare_workers_custom_domain":                         ":account_id/:id",
	"cloudflare_workers_for_platforms_dispatch_namespace":      ":account_id/:id",
	"cloudflare_workers_kv_namespace":                          ":account_id/:id",
	"cloudflare_zone":                                          ":id",
	"cloudflare_zone_dnssec":                                   ":zone_id",
	"cloudflare_zone_lockdown":                                 ":zone_id/:id",
	"cloudflare_zone_setting":                                  ":zone_id/:id",
	"cloudflare_zero_trust_access_application":                 ":identifier_type/:identifier_value/:id",
	"cloudflare_zero_trust_access_custom_page":                 ":account_id/:id",
	"cloudflare_zero_trust_access_infrastructure_target":       ":account_id/:id",
	"cloudflare_zero_trust_access_key_configuration":           ":account_id",
	"cloudflare_zero_trust_access_policy":                      ":account_id/:id",
	"cloudflare_zero_trust_access_short_lived_certificate":     ":identifier_type/:identifier_value/:id",
	"cloudflare_zero_trust_access_tag":                         ":account_id/:id",
	"cloudflare_zero_trust_device_custom_profile":              ":account_id/:id",
	"cloudflare_zero_trust_device_default_profile":             ":account_id",
	"cloudflare_zero_trust_device_managed_networks":            ":account_id/:id",
	"cloudflare_zero_trust_device_posture_integration":         ":account_id/:id",
	"cloudflare_zero_trust_device_posture_rule":                ":account_id/:id",
	"cloudflare_zero_trust_dex_test":                           ":account_id/:id",
	"cloudflare_zero_trust_dlp_predefined_profile":             ":account_id/:id",
	"cloudflare_zero_trust_dns_location":                       ":account_id/:id",
	"cloudflare_zero_trust_gateway_certificate":                ":account_id/:id",
	"cloudflare_zero_trust_gateway_policy":                     ":account_id/:id",
	"cloudflare_zero_trust_gateway_proxy_endpoint":             ":account_id/:id",
	"cloudflare_zero_trust_gateway_settings":                   ":account_id",
	"cloudflare_zero_trust_list":                               ":account_id/:id",
	"cloudflare_zero_trust_risk_scoring_integration":           ":account_id/:id",
	"cloudflare_zero_trust_tunnel_cloudflared":                 ":account_id/:id",
	"cloudflare_zero_trust_tunnel_cloudflared_route":           ":account_id/:id",
	"cloudflare_zero_trust_tunnel_cloudflared_virtual_network": ":account_id/:id",
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
		execPath, err := findOrInstallTerraform()
		if err != nil {
			log.Fatalf("Could not find or install Terraform: %v", err)
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

		providerVersionString = detectedVersion.String()
		log.WithFields(logrus.Fields{
			"version":  providerVersionString,
			"registry": registryPath,
		}).Debug("detected provider")

		resourceIDsMap := make(map[string][]string)
		var (
			jsonStructData                       []interface{}
			pathParams, endpointsWithResourceIDs []string
		)

		if strings.HasPrefix(providerVersionString, "5") {
			resources := strings.Split(resourceType, ",")
			for _, resourceType := range resources {
				if isSupportedPathParam(resources, resourceType) {
					resourceIDsMap = getResourceMappings()
					pathParams, ok = resourceIDsMap[resourceType]
					if ok && len(pathParams) == 0 {
						log.Fatalf("No resource IDs defined in Terraform for resource %s", resourceType)
					}
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

				if apiToken != "" {
					api.Options = append(api.Options, option.WithAPIToken(apiToken))
				} else {
					api.Options = append(api.Options, option.WithAPIKey(apiKey), option.WithAPIEmail(apiEmail))
				}

				if len(pathParams) > 0 {
					endpointsWithResourceIDs = replacePathParams(pathParams, endpoint, resourceType)
					jsonStructData, err = getAPIResponse(result, pathParams, endpointsWithResourceIDs...)
					if err != nil {
						log.Infof("error getting API response for resource %s: %s", resourceType, err)
						continue
					}
				} else {
					jsonStructData, err = getAPIResponse(result, pathParams, endpoint)
					if err != nil {
						log.Infof("error getting API response for resource %s: %s", resourceType, err)
						continue
					}
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
					jsonPayload, _, err := apiV0.ListAccessApplications(context.Background(), identifier, cfv0.ListAccessApplicationsParams{})
					if err != nil {
						log.Fatal(err)
					}

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
					botManagement, err := apiV0.GetBotManagement(context.Background(), identifier)
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
					jsonPayload, err := apiV0.ListPrefixes(context.Background(), accountID)
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
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

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_custom_pages":
					if accountID != "" {
						jsonPayload, err := apiV0.CustomPages(context.Background(), &cfv0.CustomPageOptions{AccountID: accountID})
						if err != nil {
							log.Fatal(err)
						}

						m, _ := json.Marshal(jsonPayload)
						err = json.Unmarshal(m, &jsonStructData)
						if err != nil {
							log.Fatal(err)
						}
					} else {
						jsonPayload, err := apiV0.CustomPages(context.Background(), &cfv0.CustomPageOptions{ZoneID: zoneID})
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
					jsonPayload, _, err := apiV0.Filters(context.Background(), identifier, cfv0.FilterListParams{})
					if err != nil {
						log.Fatal(err)
					}
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
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_custom_hostname":
					jsonPayload, _, err := apiV0.CustomHostnames(context.Background(), zoneID, 1, cfv0.CustomHostname{})
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_custom_ssl":
					jsonPayload, err := apiV0.ListSSL(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_ip_list":
					jsonPayload, err := apiV0.ListIPLists(context.Background(), accountID)
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_load_balancer":
					jsonPayload, err := apiV0.ListLoadBalancers(context.Background(), identifier, cfv0.ListLoadBalancerParams{})
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_load_balancer_pool":
					jsonPayload, err := apiV0.ListLoadBalancerPools(context.Background(), identifier, cfv0.ListLoadBalancerPoolParams{})
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_load_balancer_monitor":
					jsonPayload, err := apiV0.ListLoadBalancerMonitors(context.Background(), identifier, cfv0.ListLoadBalancerMonitorParams{})
					if err != nil {
						log.Fatal(err)
					}
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
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_origin_ca_certificate":
					jsonPayload, err := apiV0.ListOriginCACertificates(context.Background(), cfv0.ListOriginCertificatesParams{ZoneID: zoneID})
					if err != nil {
						log.Fatal(err)
					}

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

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_rate_limit":
					jsonPayload, err := apiV0.ListAllRateLimits(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_record":
					jsonPayload, _, err := apiV0.ListDNSRecords(context.Background(), identifier, cfv0.ListDNSRecordsParams{})
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_ruleset":
					jsonPayload, err := apiV0.ListRulesets(context.Background(), identifier, cfv0.ListRulesetsParams{})
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
					jsonPayload, err := apiV0.SpectrumApplications(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}

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

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_teams_location":
					jsonPayload, _, err := apiV0.TeamsLocations(context.Background(), accountID)
					if err != nil {
						log.Fatal(err)
					}

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

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_tunnel":
					log.Debug("only requesting the first 1000 active Cloudflare Tunnels due to the service not providing correct pagination responses")
					jsonPayload, _, err := apiV0.ListTunnels(
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
					jsonPayload, _, err := apiV0.ListTurnstileWidgets(context.Background(), identifier, cfv0.ListTurnstileWidgetParams{})
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
					jsonPayload, err := apiV0.ListWAFOverrides(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}

					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_waf_package":
					jsonPayload, err := apiV0.ListWAFPackages(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_waiting_room":
					jsonPayload, err := apiV0.ListWaitingRooms(context.Background(), zoneID)
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_workers_kv_namespace":
					jsonPayload, _, err := apiV0.ListWorkersKVNamespaces(context.Background(), identifier, cfv0.ListWorkersKVNamespacesParams{})
					if err != nil {
						log.Fatal(err)
					}

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

					m, _ := json.Marshal(jsonPayload.Routes)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_zone":
					jsonPayload, err := apiV0.ListZones(context.Background())
					if err != nil {
						log.Fatal(err)
					}
					m, _ := json.Marshal(jsonPayload)
					err = json.Unmarshal(m, &jsonStructData)
					if err != nil {
						log.Fatal(err)
					}
				case "cloudflare_zone_lockdown":
					jsonPayload, _, err := apiV0.ListZoneLockdowns(context.Background(), identifier, cfv0.LockdownListParams{})
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
		skippedForMissingID := 0
		for i, data := range jsonStructData {
			var id string

			if data.(map[string]interface{})["id"] == nil {
				// Not every API object names its identifier "id". Where the
				// resource's identifier lives under a different key, use it.
				id = identifierFromResponse(resourceType, data.(map[string]interface{}))

				if id == "" {
					// Falling back to the account or zone ID is only correct when
					// that IS the resource's import id -- a per-account or
					// per-zone singleton such as cloudflare_account or
					// cloudflare_email_routing_dns. When the import format
					// contains :id, the scope id is definitionally the wrong
					// value, and emitting it produces an import block that is
					// individually plausible and collectively wrong: every object
					// of the type gets the same id, and because the resource name
					// is derived from the same value, they collide too.
					//
					// Skip rather than emit a known-bad id. A missing import block
					// is visible; a wrong one silently adopts the wrong object.
					if format, known := resourceImportStringFormats[resourceType]; known && strings.Contains(format, ":id") {
						skippedForMissingID++
						log.Errorf("%s: no identifier found in the API response for entry %d "+
							"(looked for \"id\"%s). Import format is %q, which requires one, so "+
							"skipping this resource rather than emitting the account or zone ID in "+
							"its place. Please report this resource type as a bug.",
							resourceType, i, idFieldHint(resourceType), format)
						continue
					}

					if accountID != "" {
						id = accountID
					}

					if zoneID != "" {
						id = zoneID
					}
				}
			} else {
				switch data.(map[string]interface{})["id"].(type) {
				case float64:
					id = fmt.Sprintf("%d", int(data.(map[string]interface{})["id"].(float64)))
				default:
					id = data.(map[string]interface{})["id"].(string)
				}
			}
			// Name scope-addressed resources after their scope, not after the
			// object's own identifier, or every zone produces the same address.
			nameID := id
			if scoped := scopeAddressedIdentifier(resourceToEndpoint[resourceType]["get"]); scoped != "" {
				nameID = scoped
			}

			if useModernImportBlock {
				idvalue := buildRawImportAddress(resourceType, id, resourceToEndpoint[resourceType]["get"])
				imp := importBody.AppendNewBlock("import", []string{}).Body()
				imp.SetAttributeRaw("to", hclwrite.TokensForIdentifier(fmt.Sprintf("%s.%s", resourceType, fmt.Sprintf("%s_%s_%d", terraformResourceNamePrefix, hclSafeIdentifier(nameID), i))))
				imp.SetAttributeValue("id", cty.StringVal(idvalue))
				importFile.Body().AppendNewline()
			} else {
				_, _ = fmt.Fprint(cmd.OutOrStdout(), buildTerraformImportCommand(i, resourceType, id, resourceToEndpoint[resourceType]["get"]))
			}
		}

		if skippedForMissingID > 0 {
			// A scripted `cf-terraforming import > import.tf` would otherwise produce
			// a file quietly missing resources, with a zero exit status. The per-entry
			// errors above go to stderr; this makes the omission impossible to miss
			// when stderr is being skimmed rather than read.
			log.Errorf("%s: skipped %d resource(s) whose identifier could not be determined. "+
				"The generated import blocks are INCOMPLETE for this resource type.",
				resourceType, skippedForMissingID)
		}

		if useModernImportBlock {
			// don't format the output; there is a bug in hclwrite.Format that
			// splits incorrectly on certain characters. instead, manually
			// insert new lines on the block.
			_, _ = fmt.Fprint(cmd.OutOrStdout(), string(importFile.Bytes()))
		}
	}
}

// buildTerraformImportCommand takes the resourceType and resourceID in order to
// look up the resource type import string and then return a suitable composite
// value that is compatible with `terraform import`.
//
// Note: `endpoint` is only used on > v4. Otherwise, it is ignored.
func buildTerraformImportCommand(i int, resourceType, resourceID, endpoint string) string {
	resourceImportAddress := buildRawImportAddress(resourceType, resourceID, endpoint)
	return fmt.Sprintf("%s %s.%s %s\n", terraformImportCmdPrefix, resourceType, fmt.Sprintf("%s_%s_%d", terraformResourceNamePrefix, resourceID, i), resourceImportAddress)
}

// buildRawImportAddress takes the resourceType and resourceID in order to look up
// the resource type import string and then return a suitable address.
//
// Note: `endpoint` is only used on > v4. Otherwise, it is ignored.
func buildRawImportAddress(resourceType, resourceID, endpoint string) string {
	if strings.HasPrefix(providerVersionString, "5") {
		prefix := ""
		if strings.Contains(endpoint, "{accounts_or_zones}") {
			if accountID != "" {
				prefix = "accounts"
				endpoint = strings.Replace(endpoint, "/{accounts_or_zones}/{account_or_zone_id}/", "/accounts/{account_id}/", 1)
			} else {
				prefix = "zones"
				endpoint = strings.Replace(endpoint, "/{accounts_or_zones}/{account_or_zone_id}/", "/zones/{zone_id}/", 1)
			}
		}

		r, _ := regexp.Compile("({[a-z0-9_]*})")
		matches := r.FindAllString(endpoint, -1)

		if len(matches) > 0 {
			// A resource whose endpoint carries only a scope placeholder is addressed
			// by that scope and nothing else: /zones/{zone_id}/argo/tiered_caching is
			// one object per zone, and the provider imports it as "<zone_id>". Filling
			// the placeholder with the resource ID yields the object's own identifier
			// -- a setting name such as "tiered_cache_smart_topology_enable", or a
			// rule id -- which is not a zone, so every zone emits the same value and
			// the import adopts the wrong object. The prefix is dropped for the same
			// reason: the import format is the bare scope id, not "accounts/<id>".
			if len(matches) == 1 && isScopePlaceholder(matches[0]) {
				if matches[0] == "{account_id}" {
					return accountID
				}
				return zoneID
			}

			// Naive assumptions below but if we only have a single placeholder (`{}`)
			// we can replace that with the `resourceID` however, if we have more than
			// a single one, we assume it is the second match since that is our URL
			// conventions.
			//
			// Note: this will likely break on un-RESTful routes.
			if len(matches) == 1 {
				matches[0] = resourceID
			} else {
				if matches[0] == "{account_id}" {
					matches[0] = accountID
				} else if matches[0] == "{zone_id}" {
					matches[0] = zoneID
				}
				matches[1] = resourceID
			}
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
