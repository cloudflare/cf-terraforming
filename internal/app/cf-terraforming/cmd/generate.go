package cmd

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
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
		zoneID = viper.GetString("zone")
		accountID = viper.GetString("account")

		tmpDir, err := ioutil.TempDir("", "tfinstall")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		execPath, err := tfinstall.Find(context.Background(), tfinstall.LatestVersion(tmpDir, false))
		if err != nil {
			log.Fatal(err)
		}

		// Setup and configure Terraform to operate in the temporary directory where
		// the provider is already configured.
		workingDir := viper.GetString("terraform-install-path")
		log.Debugf("initialising Terraform in %s", workingDir)
		tf, err := tfexec.NewTerraform(workingDir, execPath)
		if err != nil {
			log.Fatal(err)
		}

		err = tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.LockTimeout("60s"))
		if err != nil {
			log.Fatal(err)
		}

		log.Debug("reading Terraform schema for Cloudflare provider")
		ps, err := tf.ProvidersSchema(context.Background())
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
		case "cloudflare_argo_tunnel":
			jsonPayload, err := api.ArgoTunnels(context.Background(), accountID)
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_byo_ip_prefix":
			jsonPayload, err := api.ListPrefixes(context.Background())
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
		case "cloudflare_load_balancer_pool":
			jsonPayload, err := api.ListLoadBalancerPools(context.Background())
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
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
		case "cloudflare_origin_ca_certificate":
			jsonPayload, err := api.OriginCertificates(context.Background(), cloudflare.OriginCACertificateListOptions{ZoneID: zoneID})
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_rate_limit":
			jsonPayload, _, err := api.ListRateLimits(context.Background(), zoneID, cloudflare.PaginationOptions{})
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_record":
			simpleDNSTypes := []string{"A", "AAAA", "CNAME", "TXT", "MX", "NS"}
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
					jsonStructData[i].(map[string]interface{})["value"] = jsonStructData[i].(map[string]interface{})["content"]
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
			}
		case "cloudflare_zone_lockdown":
			jsonPayload, err := api.ListZoneLockdowns(context.Background(), zoneID, 1)
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload.Result)
			m, _ := json.Marshal(jsonPayload.Result)
			json.Unmarshal(m, &jsonStructData)
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
				resourceID = fmt.Sprintf("terraform_managed_resource_%s", structData["id"].(string))
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

				// Skip unusable timestamps
				if contains([]string{"modified_on", "created_on"}, attrName) {
					continue
				}

				if attrName == "account_id" && accountID != "" {
					output += writeAttrLine(attrName, accountID, 2, false)
					continue
				}

				if attrName == "zone_id" && zoneID != "" {
					output += writeAttrLine(attrName, zoneID, 2, false)
					continue
				}

				ty := r.Block.Attributes[attrName].AttributeType
				switch {
				case ty.IsPrimitiveType():
					switch ty {
					case cty.String, cty.Bool, cty.Number:
						output += writeAttrLine(attrName, structData[attrName], 2, false)
					default:
						log.Debugf("unexpected primitive type %q", ty.FriendlyName())
					}
				case ty.IsCollectionType():
					switch {
					case ty.IsListType(), ty.IsSetType():
						output += writeAttrLine(attrName, structData[attrName], 2, false)
					case ty.IsMapType():
						output += writeAttrLine(attrName, structData[attrName], 2, false)
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

			sortedNestedBlocks := make([]string, 0, len(r.Block.NestedBlocks))
			for k := range r.Block.NestedBlocks {
				sortedNestedBlocks = append(sortedNestedBlocks, k)
			}
			sort.Strings(sortedNestedBlocks)

			// Nested blocks are used for configuration options where assignment
			// isn't required.
			for _, attrName := range sortedNestedBlocks {
				structData := jsonStructData[i].(map[string]interface{})

				if r.Block.NestedBlocks[attrName].NestingMode == "list" || r.Block.NestedBlocks[attrName].NestingMode == "set" {
					sortedInnerNestedBlock := make([]string, 0, len(r.Block.NestedBlocks[attrName].Block.Attributes))
					for k := range r.Block.NestedBlocks[attrName].Block.Attributes {
						sortedInnerNestedBlock = append(sortedInnerNestedBlock, k)
					}
					sort.Strings(sortedInnerNestedBlock)

					nestedBlockOutput := ""
					for _, nestedAttrName := range sortedInnerNestedBlock {
						ty := r.Block.NestedBlocks[attrName].Block.Attributes[nestedAttrName].AttributeType
						switch {
						case ty.IsPrimitiveType():
							switch ty {
							case cty.String, cty.Bool, cty.Number:
								if structData[attrName] != nil {
									switch structData[attrName].(type) {
									case map[string]interface{}:
										nestedBlockOutput += writeAttrLine(nestedAttrName, structData[attrName].(map[string]interface{})[nestedAttrName], 4, false)
									default:
										log.Debugf("unexpected nested primitive type %T for %s", structData[attrName], attrName)
									}
								}
							default:
								log.Debugf("unexpected primitive type %q", ty.FriendlyName())
							}
						}
					}

					if nestedBlockOutput != "" {
						output += "  " + attrName + " {\n"
						output += nestedBlockOutput
						output += "  }\n"
					}

				} else {
					log.Debugf("nested mode %q for %s not recognised", r.Block.NestedBlocks[attrName].NestingMode, attrName)
				}
			}

			output += "}\n\n"
		}

		fmt.Fprint(cmd.OutOrStdout(), output)
	}
}

// writeAttrLine outputs a line of HCL configuration with a configurable depth
// for known types.
func writeAttrLine(key string, value interface{}, depth int, usedInBlock bool) string {
	switch value.(type) {
	case map[string]interface{}:

		values := value.(map[string]interface{})

		sortedKeys := make([]string, 0, len(values))
		for k := range values {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)

		s := ""
		for _, v := range sortedKeys {
			s += writeAttrLine(v, values[v], depth+2, false)
		}

		if usedInBlock {
			if s != "" {
				return fmt.Sprintf("%s%s {\n%s%s}\n", strings.Repeat(" ", depth), key, s, strings.Repeat(" ", depth))
			}
		} else {
			if s != "" {
				return fmt.Sprintf("%s%s = {\n%s%s}\n", strings.Repeat(" ", depth), key, s, strings.Repeat(" ", depth))
			}
		}

	case []interface{}:
		var items []string
		for _, item := range value.([]interface{}) {
			items = append(items, fmt.Sprintf("%q", item.(string)))
		}

		if len(items) > 0 {
			return fmt.Sprintf("%s%s = [ %s ]\n", strings.Repeat(" ", depth), key, strings.Join(items, ", "))
		}
	case []string:
		var items []string
		for _, item := range value.([]string) {
			items = append(items, fmt.Sprintf("%q", item))
		}
		if len(items) > 0 {
			return fmt.Sprintf("%s%s = [ %s ]\n", strings.Repeat(" ", depth), key, strings.Join(items, ", "))
		}
	case string:
		if value != "" {
			return fmt.Sprintf("%s%s = %q\n", strings.Repeat(" ", depth), key, value)
		}
	case int:
		return fmt.Sprintf("%s%s = %d\n", strings.Repeat(" ", depth), key, value)
	case float64:
		return fmt.Sprintf("%s%s = %0.f\n", strings.Repeat(" ", depth), key, value)
	case bool:
		return fmt.Sprintf("%s%s = %t\n", strings.Repeat(" ", depth), key, value)
	default:
		log.Debugf("got unknown attribute configuration: key %s, value %v, value type %T", key, value, value)
		return ""
	}
	return ""
}
