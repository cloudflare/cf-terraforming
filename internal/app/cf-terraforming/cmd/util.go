package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	cloudflare "github.com/cloudflare/cloudflare-go"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zclconf/go-cty/cty"
)

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

// testDataFile slurps a local test case into memory and returns it while
// encapsulating the logic for finding it.
func testDataFile(filename string) string {
	filename = strings.TrimSuffix(filename, "/")

	dirname, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dir, err := os.Open(filepath.Join(dirname, "../../../../testdata/terraform"))
	if err != nil {
		panic(err)
	}

	fullpath := dir.Name() + "/" + filename
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		panic(fmt.Errorf("terraform testdata file does not exist at %s", fullpath))
	}

	data, _ := ioutil.ReadFile(fullpath)

	return string(data)
}

func sharedPreRun(cmd *cobra.Command, args []string) {
	accountID = viper.GetString("account")
	zoneID = viper.GetString("zone")

	if accountID != "" && zoneID != "" {
		log.Fatal("--account and --zone are mutually exclusive and cannot be used together")
	}

	if apiToken = viper.GetString("token"); apiToken == "" {
		if apiEmail = viper.GetString("email"); apiEmail == "" {
			log.Error("'email' must be set.")
		}

		if apiKey = viper.GetString("key"); apiKey == "" {
			log.Error("either -t/--token or -k/--key must be set.")
		}

		log.WithFields(logrus.Fields{
			"email":      apiEmail,
			"zone_id":    zoneID,
			"account_id": accountID,
		}).Debug("initializing cloudflare-go")

	} else {
		log.WithFields(logrus.Fields{
			"zone_id":    zoneID,
			"account_Id": accountID,
		}).Debug("initializing cloudflare-go with API Token")
	}

	var options []cloudflare.Option

	if accountID != "" {
		log.WithFields(logrus.Fields{
			"account_id": accountID,
		}).Debug("configuring Cloudflare API with account")

		// Organization ID was passed, use it to configure the API
		options = append(options, cloudflare.UsingAccount(accountID))
	}

	var apiHost string
	if apiHost = os.Getenv("CLOUDFLARE_API_HOSTNAME"); apiHost != "" {
		options = append(options, cloudflare.BaseURL("https://"+apiHost+"/client/v4"))
	}

	if verbose {
		options = append(options, cloudflare.Debug(true))
	}

	var err error

	// Don't initialise a client in CI as this messes with VCR and the ability to
	// mock out the HTTP interactions.
	if os.Getenv("CI") != "true" {
		var useToken = apiToken != ""

		if useToken {
			api, err = cloudflare.NewWithAPIToken(apiToken, options...)
		} else {
			api, err = cloudflare.New(apiKey, apiEmail, options...)
		}

		if err != nil {
			log.Fatal(err)
		}
	}
}

// sanitiseTerraformResourceName ensures that a Terraform resource name matches
// the restrictions imposed by core.
func sanitiseTerraformResourceName(s string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	return re.ReplaceAllString(s, "_")
}

// flattenAttrMap takes a list of attributes defined as a list of maps comprising of {"id": "attrId", "value": "attrValue"}
// and flattens it to a single map of {"attrId": "attrValue"}
func flattenAttrMap(l []interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	attrID := ""
	var attrVal interface{}

	for _, elem := range l {
		switch t := elem.(type) {
		case map[string]interface{}:
			if id, ok := t["id"]; ok {
				attrID = id.(string)
			} else {
				log.Debug("no 'id' in map when attempting to flattenAttrMap")
			}

			if val, ok := t["value"]; ok {
				if val == nil {
					log.Debugf("Found nil 'value' for %s attempting to flattenAttrMap, coercing to true", attrID)
					attrVal = true
				} else {
					attrVal = val
				}
			} else {
				log.Debug("no 'value' in map when attempting to flattenAttrMap")
			}

			result[attrID] = attrVal
		default:
			log.Debugf("got unknown element type %T when attempting to flattenAttrMap", elem)
		}
	}

	return result
}

// nestBlocks takes a schema and generates all of the appropriate nesting of any
// top-level blocks as well as nested lists or sets.
func nestBlocks(schemaBlock *tfjson.SchemaBlock, structData map[string]interface{}, depth int) string {
	output := ""

	// Nested blocks are used for configuration options where assignment
	// isn't required.
	sortedNestedBlocks := make([]string, 0, len(schemaBlock.NestedBlocks))
	for k := range schemaBlock.NestedBlocks {
		sortedNestedBlocks = append(sortedNestedBlocks, k)
	}
	sort.Strings(sortedNestedBlocks)

	for _, block := range sortedNestedBlocks {
		if schemaBlock.NestedBlocks[block].NestingMode == "list" || schemaBlock.NestedBlocks[block].NestingMode == "set" {
			sortedInnerAttributes := make([]string, 0, len(schemaBlock.NestedBlocks[block].Block.Attributes))

			for k := range schemaBlock.NestedBlocks[block].Block.Attributes {
				sortedInnerAttributes = append(sortedInnerAttributes, k)
			}

			sort.Strings(sortedInnerAttributes)

			nestedBlockOutput := ""

			// If the attribute we're looking at has further nesting, we'll
			// recursively call nestBlocks.
			if len(schemaBlock.NestedBlocks[block].Block.NestedBlocks) > 0 {
				if s, ok := structData[block]; ok {
					switch s.(type) {
					case map[string]interface{}:
						nestedBlockOutput += nestBlocks(schemaBlock.NestedBlocks[block].Block, s.(map[string]interface{}), depth+2)
					case []interface{}:
						for _, nestedItem := range s.([]interface{}) {
							nestedBlockOutput += nestBlocks(schemaBlock.NestedBlocks[block].Block, nestedItem.(map[string]interface{}), depth+2)
						}

					default:
						log.Debugf("unable to generate recursively nested blocks for %T", s)
					}
				}
			}

			switch attrStruct := structData[block].(type) {

			// Case for if the inner block's attributes are a map of interfaces, in
			// which case we can directly add them to the config.
			case map[string]interface{}:
				if attrStruct != nil {
					nestedBlockOutput += writeNestedBlock(sortedInnerAttributes, schemaBlock.NestedBlocks[block].Block, attrStruct, depth)
				}

				if nestedBlockOutput != "" || schemaBlock.NestedBlocks[block].MinItems > 0 {
					output += strings.Repeat(" ", depth) + block + " {\n"
					output += nestedBlockOutput
					output += strings.Repeat(" ", depth) + "}\n"
				}

			// Case for if the inner block's attributes are a list of map interfaces,
			// in which case this should be treated as a duplicating block.
			case []map[string]interface{}:
				for _, v := range attrStruct {
					repeatedBlockOutput := ""

					if attrStruct != nil {
						repeatedBlockOutput = writeNestedBlock(sortedInnerAttributes, schemaBlock.NestedBlocks[block].Block, v, depth)
					}

					// Write the block if we had data for it, or if it is a required block.
					if repeatedBlockOutput != "" || schemaBlock.NestedBlocks[block].MinItems > 0 {
						output += strings.Repeat(" ", depth) + block + " {\n"
						output += repeatedBlockOutput

						if nestedBlockOutput != "" {
							output += nestedBlockOutput
						}

						output += strings.Repeat(" ", depth) + "}\n"
					}
				}

			// Case for duplicated blocks that commonly end up as an array or list at
			// the API level.
			case []interface{}:
				for _, v := range attrStruct {
					repeatedBlockOutput := ""

					if attrStruct != nil {
						repeatedBlockOutput = writeNestedBlock(sortedInnerAttributes, schemaBlock.NestedBlocks[block].Block, v.(map[string]interface{}), depth)
					}

					// Write the block if we had data for it, or if it is a required block.
					if repeatedBlockOutput != "" || schemaBlock.NestedBlocks[block].MinItems > 0 {
						output += strings.Repeat(" ", depth) + block + " {\n"
						output += repeatedBlockOutput

						if nestedBlockOutput != "" {
							output += nestedBlockOutput
						}

						output += strings.Repeat(" ", depth) + "}\n"
					}
				}
			default:
				log.Debugf("unexpected attribute struct type %T for block %s", attrStruct, block)
			}

		} else {
			log.Debugf("nested mode %q for %s not recognised", schemaBlock.NestedBlocks[block].NestingMode, block)
		}
	}

	return output
}

func writeNestedBlock(attributes []string, schemaBlock *tfjson.SchemaBlock, attrStruct map[string]interface{}, depth int) string {
	nestedBlockOutput := ""

	for _, attrName := range attributes {
		ty := schemaBlock.Attributes[attrName].AttributeType

		switch {
		case ty.IsPrimitiveType():
			switch ty {
			case cty.String, cty.Bool, cty.Number:
				nestedBlockOutput += writeAttrLine(attrName, attrStruct[attrName], depth+2, false)
			default:
				log.Debugf("unexpected primitive type %q", ty.FriendlyName())
			}
		case ty.IsListType(), ty.IsSetType():
			nestedBlockOutput += writeAttrLine(attrName, attrStruct[attrName], depth+2, false)
		default:
			log.Debugf("unexpected nested type %T for %s", ty, attrName)
		}
	}

	return nestedBlockOutput
}
