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

// sanitiseTerraformResourceName ensures that a Terraform resource name matches the
// restrictions imposed by core.
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
				attrVal = val
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

func nestBlocks(schemaBlock *tfjson.SchemaBlock, structData map[string]interface{}, depth int) string {
	output := ""

	// Nested blocks are used for configuration options where assignment
	// isn't required.
	sortedNestedBlocks := make([]string, 0, len(schemaBlock.NestedBlocks))
	for k := range schemaBlock.NestedBlocks {
		sortedNestedBlocks = append(sortedNestedBlocks, k)
	}
	sort.Strings(sortedNestedBlocks)

	for _, attrName := range sortedNestedBlocks {
		if schemaBlock.NestedBlocks[attrName].NestingMode == "list" || schemaBlock.NestedBlocks[attrName].NestingMode == "set" {
			sortedInnerNestedBlock := make([]string, 0, len(schemaBlock.NestedBlocks[attrName].Block.Attributes))

			for k := range schemaBlock.NestedBlocks[attrName].Block.Attributes {
				sortedInnerNestedBlock = append(sortedInnerNestedBlock, k)
			}

			sort.Strings(sortedInnerNestedBlock)

			nestedBlockOutput := ""

			// If the attribute we're looking at has further nesting, we'll recursively call nestBlocks.
			if len(schemaBlock.NestedBlocks[attrName].Block.NestedBlocks) > 0 {
				if s, ok := structData[attrName]; ok {
					nestedBlockOutput += nestBlocks(schemaBlock.NestedBlocks[attrName].Block, s.(map[string]interface{}), depth+2)
				}
			}

			for _, nestedAttrName := range sortedInnerNestedBlock {
				ty := schemaBlock.NestedBlocks[attrName].Block.Attributes[nestedAttrName].AttributeType

				switch {
				case ty.IsPrimitiveType():
					switch ty {
					case cty.String, cty.Bool, cty.Number:
						if structData[attrName] != nil {
							switch structData[attrName].(type) {
							case map[string]interface{}:
								nestedBlockOutput += writeAttrLine(nestedAttrName, structData[attrName].(map[string]interface{})[nestedAttrName], depth+2, false)
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
				output += strings.Repeat(" ", depth) + attrName + " {\n"
				output += nestedBlockOutput
				output += strings.Repeat(" ", depth) + "}\n"
			}

		} else {
			log.Debugf("nested mode %q for %s not recognised", schemaBlock.NestedBlocks[attrName].NestingMode, attrName)
		}
	}

	return output
}
