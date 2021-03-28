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
	"github.com/thanhpk/randstr"
	"github.com/zclconf/go-cty/cty"

	"fmt"
	"strings"
)

const (
	errAccountIDMissing = "account_id is expected on the resource however the provided value is missing"
	errZoneIDMissing    = "zone_id is expected on the resource however the provided value is missing"
)

var (
	output       string
	resourceType string
)

func init() {
	rootCmd.AddCommand(GenerateCmd())
}

func GenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Pull resources from the Cloudflare API and generate the respective Terraform resources",
		Run: func(cmd *cobra.Command, args []string) {
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
			// the provider is already configured. Eventually, this will be '.'.
			workingDir := "/tmp"
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

			r := s.ResourceSchemas[*&resourceType]

			log.Debugf("beginning to read and build %s resources", *&resourceType)

			// Initialise `resourceCount` outside of the switch for supported resources
			// to allow it to be referenced further down in the loop that outputs the
			// newly generated resources.
			resourceCount := 0

			// Lazy approach to restrict support to known resources due to Go's type
			// restrictions and the need to explicitly map out the structs.
			var jsonStructData []interface{}
			switch *&resourceType {
			case "cloudflare_access_identity_provider":
				if *&accountID != "" {
					jsonPayload, err := api.AccessIdentityProviders(context.Background(), *&accountID)
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					json.Unmarshal(m, &jsonStructData)
				} else {
					jsonPayload, err := api.ZoneLevelAccessIdentityProviders(context.Background(), *&zoneName)
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					json.Unmarshal(m, &jsonStructData)
				}
			case "cloudflare_access_service_token":
				if *&accountID != "" {
					jsonPayload, _, err := api.AccessServiceTokens(context.Background(), *&accountID)
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					json.Unmarshal(m, &jsonStructData)
				} else {
					jsonPayload, _, err := api.ZoneLevelAccessServiceTokens(context.Background(), *&zoneName)
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					json.Unmarshal(m, &jsonStructData)
				}
			case "cloudflare_access_mutual_tls_certificate":
				jsonPayload, err := api.AccessMutualTLSCertificates(context.Background(), *&accountID)
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			case "cloudflare_access_rule":
				if *&accountID != "" {
					jsonPayload, err := api.ListAccountAccessRules(context.Background(), *&accountID, cloudflare.AccessRule{}, 1)
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload.Result)
					m, _ := json.Marshal(jsonPayload.Result)
					json.Unmarshal(m, &jsonStructData)
				} else {
					jsonPayload, err := api.ListZoneAccessRules(context.Background(), *&zoneName, cloudflare.AccessRule{}, 1)
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload.Result)
					m, _ := json.Marshal(jsonPayload.Result)
					json.Unmarshal(m, &jsonStructData)
				}
			case "cloudflare_account_member":
				jsonPayload, _, err := api.AccountMembers(context.Background(), *&accountID, cloudflare.PaginationOptions{})
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
				jsonPayload, err := api.ArgoTunnels(context.Background(), *&accountID)
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
			case "cloudflare_certificate_pack":
				jsonPayload, err := api.ListCertificatePacks(context.Background(), *&zoneName)
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			case "cloudflare_custom_hostname_fallback_origin":
				jsonPayload, _, err := api.CustomHostnames(context.Background(), *&zoneName, 1, cloudflare.CustomHostname{})
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			case "cloudflare_custom_pages":
				if *&accountID != "" {
					jsonPayload, err := api.CustomPages(context.Background(), &cloudflare.CustomPageOptions{AccountID: *&accountID})
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					json.Unmarshal(m, &jsonStructData)
				} else {
					jsonPayload, err := api.CustomPages(context.Background(), &cloudflare.CustomPageOptions{ZoneID: *&zoneName})
					if err != nil {
						log.Fatal(err)
					}

					resourceCount = len(jsonPayload)
					m, _ := json.Marshal(jsonPayload)
					json.Unmarshal(m, &jsonStructData)
				}
			case "cloudflare_filter":
				jsonPayload, err := api.Filters(context.Background(), *&zoneName, cloudflare.PaginationOptions{})
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			case "cloudflare_firewall_rule":
				jsonPayload, err := api.FirewallRules(context.Background(), *&zoneName, cloudflare.PaginationOptions{})
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
				jsonPayload, _, err := api.CustomHostnames(context.Background(), *&zoneName, 1, cloudflare.CustomHostname{})
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			case "cloudflare_ip_list":
				jsonPayload, err := api.ListIPLists(context.Background())
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			case "cloudflare_record":
				simpleDNSTypes := []string{"A", "AAAA", "CNAME", "TXT", "MX", "NS"}
				jsonPayload, err := api.DNSRecords(context.Background(), *&zoneName, cloudflare.DNSRecord{})
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)

				for i := 0; i < resourceCount; i++ {
					// We only want to remap the "value" to the "content" value for simple
					// DNS types as the aggregate types use `data` for the structure.
					if contains(simpleDNSTypes, jsonStructData[i].(map[string]interface{})["type"].(string)) {
						jsonStructData[i].(map[string]interface{})["value"] = jsonStructData[i].(map[string]interface{})["content"]
					}
				}
			case "cloudflare_waf_package":
				jsonPayload, err := api.ListWAFPackages(context.Background(), *&zoneName)
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			case "cloudflare_worker_route":
				jsonPayload, err := api.ListWorkerRoutes(context.Background(), *&zoneName)
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

				// remap "zone" to the "name" value.
				for i := 0; i < resourceCount; i++ {
					jsonStructData[i].(map[string]interface{})["zone"] = jsonStructData[i].(map[string]interface{})["name"]
				}
			default:
				fmt.Fprintf(cmd.OutOrStdout(), "%q is not yet supported for automatic generation", *&resourceType)
				return
			}

			// If we don't have any resources to generate, just bail out early.
			if resourceCount == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "no resources found to generate. Exiting...")
				return
			}

			for i := 0; i < resourceCount; i++ {
				resourceID := ""
				if os.Getenv("USE_STATIC_RESOURCE_IDS") == "true" {
					resourceID = "terraform_managed_resource"
				} else {
					resourceID = fmt.Sprintf("terraform_managed_resource_%s", randstr.Hex(5))
				}
				output += fmt.Sprintf(`resource "%s" "%s" {`+"\n", *&resourceType, resourceID)

				// Block attributes are for any attributes where assignment is involved.
				for attrName, attrConfig := range r.Block.Attributes {
					// Don't bother outputting the ID for the resource as that is only for
					// internal use (such as importing state).
					if attrName == "id" {
						continue
					}

					structData := jsonStructData[i].(map[string]interface{})

					if attrName == "account_id" && *&accountID == "" {
						if *&accountID == "" {
							log.Fatal(errAccountIDMissing)
						} else {
							output += writeAttrLine(attrName, *&accountID, 2, false)
							continue
						}
					}

					if attrName == "zone_id" {
						if *&zoneName == "" {
							log.Fatal(errZoneIDMissing)
						} else {
							output += writeAttrLine(attrName, *&zoneName, 2, false)
							continue
						}
					}

					ty := attrConfig.AttributeType
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

				// Nested blocks are used for configuration options where assignment
				// isn't required.
				for attrName, attrConfig := range r.Block.NestedBlocks {
					structData := jsonStructData[i].(map[string]interface{})

					if attrConfig.NestingMode == "list" {
						output += "  " + attrName + " {\n"

						for nestedAttrName, attrConfig := range attrConfig.Block.Attributes {
							ty := attrConfig.AttributeType
							switch {
							case ty.IsPrimitiveType():
								switch ty {
								case cty.String, cty.Bool, cty.Number:
									output += writeAttrLine(nestedAttrName, structData[attrName].(map[string]interface{})[nestedAttrName], 4, false)
								default:
									log.Debugf("unexpected primitive type %q", ty.FriendlyName())
								}
							}
						}

						output += "  }\n"
					} else {
						log.Debugf("nested mode %q for %s not recongised", attrConfig.NestingMode, attrName)
					}
				}
				output += "}\n\n"
			}

			fmt.Fprintf(cmd.OutOrStdout(), output)
		},
	}

	cmd.PersistentFlags().StringVar(&resourceType, "resource-type", "", "Which resource you wish to generate")
	return cmd
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
			return fmt.Sprintf("%s%s {\n%s%s}\n", strings.Repeat(" ", depth), key, s, strings.Repeat(" ", depth))
		} else {
			return fmt.Sprintf("%s%s = {\n%s%s}\n", strings.Repeat(" ", depth), key, s, strings.Repeat(" ", depth))
		}

	case []interface{}:
		var items []string
		for _, item := range value.([]interface{}) {
			items = append(items, fmt.Sprintf("%q", item.(string)))
		}
		return fmt.Sprintf("%s%s = [ %s ]\n", strings.Repeat(" ", depth), key, strings.Join(items, ", "))
	case []string:
		var items []string
		for _, item := range value.([]string) {
			items = append(items, fmt.Sprintf("%q", item))
		}
		return fmt.Sprintf("%s%s = [ %s ]\n", strings.Repeat(" ", depth), key, strings.Join(items, ", "))
	case string:
		return fmt.Sprintf("%s%s = %q\n", strings.Repeat(" ", depth), key, value)
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
}
