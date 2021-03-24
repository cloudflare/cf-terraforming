package cmd

import (
	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/spf13/cobra"
	"github.com/thanhpk/randstr"
	"github.com/zclconf/go-cty/cty"

	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
)

const (
	errAccountIDMissing = "account_id is expected on the resource however the provided value is missing"
	errZoneIDMissing    = "zone_id is expected on the resource however the provided value is missing"
)

var (
	output       string
	resourceType string

	// schemaToAPIMapping contains an override mapping for the API <> schema
	// mismatches. The top level map key is the resource name while the inner map
	// is the API value that you wish to map to the schema.
	schemaToAPIMapping = map[string]map[string]string{
		"cloudflare_record": map[string]string{
			"value": "content", // remap "value" from the API to "content" in the schema
		},
	}
)

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().StringVarP(&resourceType, "resource-type", "r", "", "Which resource you wish to generate")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Pull resources from the Cloudflare API and generate the respective Terraform resources",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("attempting to generating %q resources", *&resourceType)

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
		case "cloudflare_access_service_token":
			if *&accountID != "" {
				jsonPayload, _, err := api.AccessServiceTokens(*&accountID)
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			} else {
				jsonPayload, _, err := api.ZoneLevelAccessServiceTokens(*&zoneName)
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			}
		case "cloudflare_argo_tunnel":
			jsonPayload, err := api.ArgoTunnels(context.Background(), *&accountID)
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		// need to work out how to zone/account differences
		// case "cloudflare_authenticated_origin_pulls":
		case "cloudflare_byo_ip_prefix":
			jsonPayload, err := api.ListPrefixes(context.Background())
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_custom_hostname_fallback_origin":
			jsonPayload, _, err := api.CustomHostnames(*&zoneName, 1, cloudflare.CustomHostname{})
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_custom_pages":
			if *&accountID != "" {
				jsonPayload, err := api.CustomPages(&cloudflare.CustomPageOptions{AccountID: *&accountID})
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			} else {
				jsonPayload, err := api.CustomPages(&cloudflare.CustomPageOptions{ZoneID: *&zoneName})
				if err != nil {
					log.Fatal(err)
				}

				resourceCount = len(jsonPayload)
				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			}
		case "cloudflare_filter":
			jsonPayload, err := api.Filters(*&zoneName, cloudflare.PaginationOptions{})
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_firewall_rule":
			jsonPayload, err := api.FirewallRules(*&zoneName, cloudflare.PaginationOptions{})
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)

			// remap Filter.ID to `filter_id` on the JSON payload
			for i := 0; i < resourceCount; i++ {
				jsonStructData[i].(map[string]interface{})["filter_id"] = jsonStructData[i].(map[string]interface{})["filter"].(map[string]interface{})["id"]
			}
		case "cloudflare_record":
			jsonPayload, err := api.DNSRecords(*&zoneName, cloudflare.DNSRecord{})
			if err != nil {
				log.Fatal(err)
			}

			resourceCount = len(jsonPayload)
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		// case "cloudflare_filter":
		// 	json.Unmarshal(jsonPayload, &jsonStructData)
		// case "cloudflare_certificate_pack":
		// 	jsonStructData = cloudflare.CertificatePack{}
		// 	json.Unmarshal(jsonPayload, &jsonStructData)
		// case "cloudflare_argo_tunnel":
		// 	jsonStructData = cloudflare.ArgoTunnel{}
		// 	json.Unmarshal(jsonPayload, &jsonStructData)
		// case "cloudflare_authenticated_origin_pulls":
		// 	jsonStructData = cloudflare.AuthenticatedOriginPulls{}
		// 	json.Unmarshal(jsonPayload, &jsonStructData)
		default:
			log.Fatalf("%q is not yet supported for automatic generation", *&resourceType)
		}

		// If we don't have any resources to generate, just bail out early.
		if resourceCount == 0 {
			log.Info("no resources found to generate. Exiting...")
			os.Exit(0)
		}

		for i := 0; i < resourceCount; i++ {
			output += fmt.Sprintf(`resource "%s" "terraform_managed_resource_%s" {`+"\n", *&resourceType, randstr.Hex(5))

			for attrName, attrConfig := range r.Block.Attributes {
				// Don't bother outputting the ID for the resource as that is only for
				// internal use (such as importing state).
				if attrName == "id" {
					continue
				}

				if val, ok := schemaToAPIMapping[*&resourceType][attrName]; ok {
					attrName = val
				}

				structData := jsonStructData[i].(map[string]interface{})

				if attrName == "account_id" && *&accountID == "" {
					if *&accountID == "" {
						log.Fatalf(errAccountIDMissing)
					} else {
						output += writeAttrLine(attrName, *&accountID, 2)
						continue
					}
				}

				if attrName == "zone_id" {
					if *&zoneName == "" {
						log.Fatalf(errZoneIDMissing)
					} else {
						output += writeAttrLine(attrName, *&zoneName, 2)
						continue
					}
				}

				ty := attrConfig.AttributeType
				switch {
				case ty.IsPrimitiveType():
					switch ty {
					case cty.String, cty.Bool, cty.Number:
						output += writeAttrLine(attrName, structData[attrName], 2)
					default:
						log.Warnf("unexpected primitive type %q", ty.FriendlyName())
					}
				case ty.IsCollectionType():
					switch {
					case ty.IsListType(), ty.IsSetType():
						output += writeAttrLine(attrName, structData[attrName], 2)
					case ty.IsMapType():
						fmt.Printf("map found. attrName %s\n", attrName)
					default:
						log.Warnf("unexpected collection type %q", ty.FriendlyName())
					}
				case ty.IsTupleType():
					fmt.Printf("tuple found. attrName %s\n", attrName)
				case ty.IsObjectType():
					fmt.Printf("object found. attrName %s\n", attrName)
				default:
					log.Warnf("attribute %q (attribute type of %q) has not been generated", attrName, ty.FriendlyName())
				}
			}

			output += "}\n\n"
		}

		fmt.Println(output)
	},
}

// writeAttrLine outputs a line of HCL configuration with a configurable depth
// for known types.
func writeAttrLine(key string, value interface{}, depth int) string {
	switch value.(type) {
	case []interface{}:
		var items []string
		for _, item := range value.([]interface{}) {
			items = append(items, fmt.Sprintf("%q", item.(string)))
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
