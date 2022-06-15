package cmd

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/google/uuid"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zclconf/go-cty/cty"

	"fmt"
	"strings"
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

		tmpDir, err := ioutil.TempDir("", "tfinstall")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		installer := &releases.ExactVersion{
			Product: product.Terraform,
			Version: version.Must(version.NewVersion("1.0.6")),
		}

		execPath, err := installer.Install(context.Background())
		if err != nil {
			log.Fatalf("error installing Terraform: %s", err)
		}

		// Setup and configure Terraform to operate in the temporary directory where
		// the provider is already configured.
		workingDir := viper.GetString("terraform-install-path")
		log.Debugf("initializing Terraform in %s", workingDir)
		tf, err := tfexec.NewTerraform(workingDir, execPath)
		if err != nil {
			log.Fatal(err)
		}

		err = tf.Init(context.Background(), tfexec.Upgrade(true))
		if err != nil {
			log.Fatal(err)
		}

		log.Debug("reading Terraform schema for Cloudflare provider")
		ps, err := tf.ProvidersSchema(context.Background())
		if err != nil {
			log.Fatal("failed to read provider schema", err)
		}
		s := ps.Schemas["registry.terraform.io/cloudflare/cloudflare"]
		if s == nil {
			log.Fatal("failed to detect provider installation")
		}

		r := s.ResourceSchemas[resourceType]
		log.Debugf("beginning to read and build %s resources", resourceType)

		// Initialise `resourceCount` outside of the switch for supported resources
		// to allow it to be referenced further down in the loop that outputs the
		// newly generated resources.
		resourceCount := 0

		// Lazy approach to restrict support to known resources due to Go's type
		// restrictions and the need to explicitly map out the structs.
		var jsonStructData []interface{}

		switch resourceType {
		case "cloudflare_access_application":
			if accountID != "" {
				jsonPayload, _, err := api.AccessApplications(context.Background(), accountID, cloudflare.PaginationOptions{})
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			} else {
				jsonPayload, _, err := api.ZoneLevelAccessApplications(context.Background(), zoneID, cloudflare.PaginationOptions{})
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			}
		case "cloudflare_access_group":
			if accountID != "" {
				jsonPayload, _, err := api.AccessGroups(context.Background(), accountID, cloudflare.PaginationOptions{})
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			} else {
				jsonPayload, _, err := api.ZoneLevelAccessGroups(context.Background(), zoneID, cloudflare.PaginationOptions{})
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			}
		case "cloudflare_access_identity_provider":
			if accountID != "" {
				jsonPayload, err := api.AccessIdentityProviders(context.Background(), accountID)
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			} else {
				jsonPayload, err := api.ZoneLevelAccessIdentityProviders(context.Background(), zoneID)
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			}
		case "cloudflare_access_service_token":
			if accountID != "" {
				jsonPayload, _, err := api.AccessServiceTokens(context.Background(), accountID)
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			} else {
				jsonPayload, _, err := api.ZoneLevelAccessServiceTokens(context.Background(), zoneID)
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			}
		case "cloudflare_access_mutual_tls_certificate":
			jsonPayload, err := api.AccessMutualTLSCertificates(context.Background(), accountID)
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_access_rule":
			if accountID != "" {
				jsonPayload, err := api.ListAccountAccessRules(context.Background(), accountID, cloudflare.AccessRule{}, 1)
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload.Result)
				m, _ := json.Marshal(jsonPayload.Result)
				json.Unmarshal(m, &jsonStructData)
			} else {
				jsonPayload, err := api.ListZoneAccessRules(context.Background(), zoneID, cloudflare.AccessRule{}, 1)
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload.Result)
				m, _ := json.Marshal(jsonPayload.Result)
				json.Unmarshal(m, &jsonStructData)
			}
		case "cloudflare_account_member":
			jsonPayload, _, err := api.AccountMembers(context.Background(), accountID, cloudflare.PaginationOptions{})
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)

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
			jsonPayload := []cloudflare.ArgoFeatureSetting{}

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
			json.Unmarshal(m, &jsonStructData)

			for _, b := range jsonStructData {
				key := b.(map[string]interface{})["id"].(string)
				jsonStructData[0].(map[string]interface{})[key] = jsonStructData[0].(map[string]interface{})["value"]
			}
		case "cloudflare_argo_tunnel":
			jsonPayload, err := api.Tunnels(context.Background(), cloudflare.TunnelListParams{
				AccountID: accountID,
				IsDeleted: cloudflare.BoolPtr(false),
			})

			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)

			for i := 0; i < resourceCount; i++ {
				secret, err := api.TunnelToken(context.Background(), cloudflare.TunnelTokenParams{
					AccountID: accountID,
					ID:        jsonStructData[i].(map[string]interface{})["id"].(string),
				})
				if err != nil {
					log.Fatal(err)
				}
				jsonStructData[i].(map[string]interface{})["secret"] = secret
			}
		case "cloudflare_byo_ip_prefix":
			jsonPayload, err := api.ListPrefixes(context.Background(), accountID)
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)

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

			var customerManagedCertificates []cloudflare.CertificatePack
			for _, r := range jsonPayload {
				if r.Type != "universal" {
					customerManagedCertificates = append(customerManagedCertificates, r)
				}
			}
			jsonPayload = customerManagedCertificates

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_custom_pages":
			if accountID != "" {
				acc := cloudflare.CustomPageOptions{AccountID: accountID}
				jsonPayload, err := api.CustomPages(context.Background(), &acc)
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			} else {
				zo := cloudflare.CustomPageOptions{ZoneID: zoneID}
				jsonPayload, err := api.CustomPages(context.Background(), &zo)
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			}

			// remap ID to the "type" field
			for i := 0; i < resourceCount; i++ {
				jsonStructData[i].(map[string]interface{})["type"] = jsonStructData[i].(map[string]interface{})["id"]
			}
		case "cloudflare_custom_hostname_fallback_origin":
			var jsonPayload []cloudflare.CustomHostnameFallbackOrigin
			apiCall, err := api.CustomHostnameFallbackOrigin(context.Background(), zoneID)
			if err != nil {
				log.Fatal(err)
			}

			if apiCall.Origin != "" {
				resourceCount = 1
				jsonPayload = append(jsonPayload, apiCall)
			}

			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)

			for i := 0; i < resourceCount; i++ {
				jsonStructData[i].(map[string]interface{})["id"] = sanitiseTerraformResourceName(jsonStructData[i].(map[string]interface{})["origin"].(string))
				jsonStructData[i].(map[string]interface{})["status"] = nil
			}
		case "cloudflare_filter":
			jsonPayload, err := api.Filters(context.Background(), zoneID, cloudflare.PaginationOptions{})
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_firewall_rule":
			jsonPayload, err := api.FirewallRules(context.Background(), zoneID, cloudflare.PaginationOptions{})
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)

			// remap Filter.ID to `filter_id` on the JSON payloads.
			for i := 0; i < resourceCount; i++ {
				jsonStructData[i].(map[string]interface{})["filter_id"] = jsonStructData[i].(map[string]interface{})["filter"].(map[string]interface{})["id"]
			}
		case "cloudflare_custom_hostname":
			jsonPayload, _, err := api.CustomHostnames(context.Background(), zoneID, 1, cloudflare.CustomHostname{})
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)

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
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_healthcheck":
			jsonPayload, err := api.Healthchecks(context.Background(), zoneID)
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_load_balancer":
			jsonPayload, err := api.ListLoadBalancers(context.Background(), zoneID)
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)

			for i := 0; i < resourceCount; i++ {
				jsonStructData[i].(map[string]interface{})["default_pool_ids"] = jsonStructData[i].(map[string]interface{})["default_pools"]
				jsonStructData[i].(map[string]interface{})["fallback_pool_id"] = jsonStructData[i].(map[string]interface{})["fallback_pool"]
			}
		case "cloudflare_load_balancer_pool":
			jsonPayload, err := api.ListLoadBalancerPools(context.Background())
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)

			for i := 0; i < resourceCount; i++ {
				for originCounter := range jsonStructData[i].(map[string]interface{})["origins"].([]interface{}) {
					if jsonStructData[i].(map[string]interface{})["origins"].([]interface{})[originCounter].(map[string]interface{})["header"] != nil {
						jsonStructData[i].(map[string]interface{})["origins"].([]interface{})[originCounter].(map[string]interface{})["header"].(map[string]interface{})["header"] = "Host"
						jsonStructData[i].(map[string]interface{})["origins"].([]interface{})[originCounter].(map[string]interface{})["header"].(map[string]interface{})["values"] = jsonStructData[i].(map[string]interface{})["origins"].([]interface{})[originCounter].(map[string]interface{})["header"].(map[string]interface{})["Host"]
					}
				}
			}
		case "cloudflare_load_balancer_monitor":
			jsonPayload, err := api.ListLoadBalancerMonitors(context.Background())
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_logpush_job":
			jsonPayload, err := api.LogpushJobs(context.Background(), zoneID)
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)

			for i := 0; i < resourceCount; i++ {
				// Workaround for LogpushJob.Filter being empty with a custom
				// marshaler and returning `{"where":{}}` as the "empty" value.
				if jsonStructData[i].(map[string]interface{})["filter"] == `{"where":{}}` {
					jsonStructData[i].(map[string]interface{})["filter"] = nil
				}
			}

		case "cloudflare_origin_ca_certificate":
			jsonPayload, err := api.OriginCertificates(context.Background(), cloudflare.OriginCACertificateListOptions{ZoneID: zoneID})
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_page_rule":
			jsonPayload, err := api.ListPageRules(context.Background(), zoneID)
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)

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
						jsonStructData[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_key_fields"].(map[string]interface{})["query_string"].(map[string]interface{})["include"] = []interface{}{}
						jsonStructData[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_key_fields"].(map[string]interface{})["query_string"].(map[string]interface{})["ignore"] = false
					}
					if s, sok := c["query_string"].(map[string]interface{})["exclude"].(string); sok && s == "*" {
						jsonStructData[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_key_fields"].(map[string]interface{})["query_string"].(map[string]interface{})["exclude"] = []interface{}{}
						jsonStructData[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_key_fields"].(map[string]interface{})["query_string"].(map[string]interface{})["ignore"] = true
					}
				}
			}
		case "cloudflare_rate_limit":
			jsonPayload, _, err := api.ListRateLimits(context.Background(), zoneID, cloudflare.PaginationOptions{})
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)

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
			simpleDNSTypes := []string{"A", "AAAA", "CNAME", "TXT", "MX", "NS", "PTR"}
			jsonPayload, err := api.DNSRecords(context.Background(), zoneID, cloudflare.DNSRecord{})
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)

			for i := 0; i < resourceCount; i++ {
				// Drop the proxiable values as they are not usable
				jsonStructData[i].(map[string]interface{})["proxiable"] = nil

				if jsonStructData[i].(map[string]interface{})["name"].(string) != jsonStructData[i].(map[string]interface{})["zone_name"].(string) {
					jsonStructData[i].(map[string]interface{})["name"] = strings.ReplaceAll(jsonStructData[i].(map[string]interface{})["name"].(string), "."+jsonStructData[i].(map[string]interface{})["zone_name"].(string), "")
				}

				// We only want to remap the "value" to the "content" value for simple
				// DNS types as the aggregate types use `data` for the structure.
				if contains(simpleDNSTypes, jsonStructData[i].(map[string]interface{})["type"].(string)) {
					switch jsonStructData[i].(map[string]interface{})["type"].(string) {
					// Edge case : when TXT record contains SPF macro that contains % then escape it with an extra %
					case "TXT":
						jsonStructData[i].(map[string]interface{})["value"] = strings.Replace(jsonStructData[i].(map[string]interface{})["content"].(string), "%", "%%", -1)
					default:
						jsonStructData[i].(map[string]interface{})["value"] = jsonStructData[i].(map[string]interface{})["content"]
					}

				}
			}
		case "cloudflare_ruleset":
			if accountID != "" {
				jsonPayload, err := api.ListAccountRulesets(context.Background(), accountID)
				if err != nil {
					log.Fatal(err)
				}

				var nonManagedRules []cloudflare.Ruleset

				// A little annoying but makes more sense doing it this way. Only append
				// the non-managed rules to the usable nonManagedRules variable instead
				// of attempting to delete from an existing slice and just reassign.
				for _, r := range jsonPayload {
					if r.Kind != string(cloudflare.RulesetKindManaged) {
						nonManagedRules = append(nonManagedRules, r)
					}
				}
				jsonPayload = nonManagedRules

				for i, rule := range nonManagedRules {
					ruleset, _ := api.GetAccountRuleset(context.Background(), accountID, rule.ID)
					jsonPayload[i].Rules = ruleset.Rules
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			} else {
				jsonPayload, err := api.ListZoneRulesets(context.Background(), zoneID)
				if err != nil {
					log.Fatal(err)
				}

				var nonManagedRules []cloudflare.Ruleset

				// A little annoying but makes more sense doing it this way. Only append
				// the non-managed rules to the usable nonManagedRules variable instead
				// of attempting to delete from an existing slice and just reassign.
				for _, r := range jsonPayload {
					if r.Kind != string(cloudflare.RulesetKindManaged) {
						nonManagedRules = append(nonManagedRules, r)
					}
				}
				jsonPayload = nonManagedRules
				ruleHeaders := map[string][]map[string]interface{}{}
				for i, rule := range nonManagedRules {

					ruleset, _ := api.GetZoneRuleset(context.Background(), zoneID, rule.ID)
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

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)

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
			}

			for i := 0; i < resourceCount; i++ {
				if jsonStructData[i].(map[string]interface{})["rules"] != nil {
					for ruleCounter := range jsonStructData[i].(map[string]interface{})["rules"].([]interface{}) {
						if jsonStructData[i].(map[string]interface{})["rules"].([]interface{})[ruleCounter].(map[string]interface{})["action_parameters"] != nil {
							if jsonStructData[i].(map[string]interface{})["rules"].([]interface{})[ruleCounter].(map[string]interface{})["action_parameters"].(map[string]interface{})["overrides"] != nil {
								if jsonStructData[i].(map[string]interface{})["rules"].([]interface{})[ruleCounter].(map[string]interface{})["action_parameters"].(map[string]interface{})["overrides"].(map[string]interface{})["enabled"] == true {
									jsonStructData[i].(map[string]interface{})["rules"].([]interface{})[ruleCounter].(map[string]interface{})["action_parameters"].(map[string]interface{})["overrides"].(map[string]interface{})["status"] = "enabled"
								}

								if jsonStructData[i].(map[string]interface{})["rules"].([]interface{})[ruleCounter].(map[string]interface{})["action_parameters"].(map[string]interface{})["overrides"].(map[string]interface{})["enabled"] == false {
									jsonStructData[i].(map[string]interface{})["rules"].([]interface{})[ruleCounter].(map[string]interface{})["action_parameters"].(map[string]interface{})["overrides"].(map[string]interface{})["status"] = "disabled"
								}

								jsonStructData[i].(map[string]interface{})["rules"].([]interface{})[ruleCounter].(map[string]interface{})["action_parameters"].(map[string]interface{})["overrides"].(map[string]interface{})["enabled"] = nil
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
			json.Unmarshal(m, &jsonStructData)

			for i := 0; i < resourceCount; i++ {
				if jsonStructData[i].(map[string]interface{})["edge_ips"] != nil {
					jsonStructData[i].(map[string]interface{})["edge_ips"] = jsonStructData[i].(map[string]interface{})["edge_ips"].(map[string]interface{})["ips"]
				}
			}
		case "cloudflare_waf_override":
			jsonPayload, err := api.ListWAFOverrides(context.Background(), zoneID)
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_waf_package":
			jsonPayload, err := api.ListWAFPackages(context.Background(), zoneID)
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_waiting_room":
			jsonPayload, err := api.ListWaitingRooms(context.Background(), zoneID)
			if err != nil {
				log.Fatal(err)
			}
			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_workers_kv_namespace":
			jsonPayload, err := api.ListWorkersKVNamespaces(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_worker_route":
			jsonPayload, err := api.ListWorkerRoutes(context.Background(), zoneID)
			if err != nil {
				log.Fatal(err)
			}
			resourceCount = len(jsonPayload.Routes)
			m, _ := json.Marshal(jsonPayload.Routes)
			json.Unmarshal(m, &jsonStructData)

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
			json.Unmarshal(m, &jsonStructData)

			// - remap "zone" to the "name" value
			// - remap "plan" to "legacy_id" value
			// - drop meta and name_servers
			for i := 0; i < resourceCount; i++ {
				jsonStructData[i].(map[string]interface{})["zone"] = jsonStructData[i].(map[string]interface{})["name"]
				jsonStructData[i].(map[string]interface{})["plan"] = jsonStructData[i].(map[string]interface{})["plan"].(map[string]interface{})["legacy_id"].(string)
				jsonStructData[i].(map[string]interface{})["meta"] = nil
				jsonStructData[i].(map[string]interface{})["name_servers"] = nil
				jsonStructData[i].(map[string]interface{})["status"] = nil
			}
		case "cloudflare_zone_lockdown":
			jsonPayload, err := api.ListZoneLockdowns(context.Background(), zoneID, 1)
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload.Result)
			m, _ := json.Marshal(jsonPayload.Result)
			json.Unmarshal(m, &jsonStructData)

		case "cloudflare_zone_settings_override":
			jsonPayload, err := api.ZoneSettings(context.Background(), zoneID)
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = 1
			m, _ := json.Marshal(jsonPayload.Result)
			json.Unmarshal(m, &jsonStructData)

			zoneSettingsStruct := make(map[string]interface{})
			for _, data := range jsonStructData {
				keyName := data.(map[string]interface{})["id"].(string)
				value := data.(map[string]interface{})["value"].(interface{})
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
		default:
			fmt.Fprintf(cmd.OutOrStdout(), "%q is not yet supported for automatic generation", resourceType)
			return
		}
		// If we don't have any resources to generate, just bail out early.
		if resourceCount == 0 {
			fmt.Fprint(cmd.OutOrStdout(), "no resources found to generate. Exiting...")
			return
		}

		output := ""

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

			output += fmt.Sprintf(`resource "%s" "%s" {`+"\n", resourceType, resourceID)

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
					output += writeAttrLine(attrName, accountID, false)
					continue
				}

				if attrName == "zone_id" && zoneID != "" && accountID == "" {
					output += writeAttrLine(attrName, zoneID, false)
					continue
				}

				ty := r.Block.Attributes[attrName].AttributeType
				switch {
				case ty.IsPrimitiveType():
					switch ty {
					case cty.String, cty.Bool, cty.Number:
						output += writeAttrLine(attrName, structData[attrName], false)
					default:
						log.Debugf("unexpected primitive type %q", ty.FriendlyName())
					}
				case ty.IsCollectionType():
					switch {
					case ty.IsListType(), ty.IsSetType(), ty.IsMapType():
						output += writeAttrLine(attrName, structData[attrName], false)
					default:
						log.Debugf("unexpected collection type %q", ty.FriendlyName())
					}
				case ty.IsTupleType():
					fmt.Printf("tuple found. attrName %s\n", attrName)
				case ty.IsObjectType():
					fmt.Printf("object found. attrName %s\n", attrName)
				default:
					log.Debugf("attribute %q (attribute type of %q) has not been generated", attrName, ty.FriendlyName())
				}
			}

			output += nestBlocks(r.Block, jsonStructData[i].(map[string]interface{}), uuid.New().String(), map[string][]string{})
			output += "}\n\n"
		}

		output, err = tf.FormatString(context.Background(), output)
		if err != nil {
			log.Fatalf("failed to format output: %s", err)
		}

		fmt.Fprint(cmd.OutOrStdout(), output)

	}
}
