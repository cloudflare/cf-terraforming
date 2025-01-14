package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"

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
			tmpDir, err := ioutil.TempDir("", "tfinstall")
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

			if resourceToEndpoint[resourceType] == "" {
				log.Debugf("did not find API endpoint for %q. skipping...", resourceType)
				continue
			}

			// Initialise `resourceCount` outside of the switch for supported resources
			// to allow it to be referenced further down in the loop that outputs the
			// newly generated resources.
			resourceCount := 0

			var (
				jsonStructData []interface{}
				result         *http.Response
			)

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
			json.Unmarshal([]byte(value.String()), &jsonStructData)

			resourceCount = len(jsonStructData)
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
