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
	"github.com/hashicorp/hcl/v2/hclwrite"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zclconf/go-cty/cty"
)

var hasNumber = regexp.MustCompile("[0-9]+").MatchString

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func executeCommandC(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	_, err = root.ExecuteC()

	return buf.String(), err
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

	fullpath := dir.Name() + "/" + filename + "/test.tf"
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		panic(fmt.Errorf("terraform testdata file does not exist at %s", fullpath))
	}

	data, _ := ioutil.ReadFile(fullpath)

	return string(data)
}

func sharedPreRun(cmd *cobra.Command, args []string) {
	accountID = viper.GetString("account")
	zoneID = viper.GetString("zone")
	hostname = viper.GetString("hostname")

	if accountID != "" && zoneID != "" {
		log.Fatal("--account and --zone are mutually exclusive, support for both is deprecated")
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

	if hostname != "" {
		options = append(options, cloudflare.BaseURL("https://"+hostname+"/client/v4"))
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
// and flattens it to a single map of {"attrId": "attrValue"}.
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

func processBlocks(schemaBlock *tfjson.SchemaBlock, structData map[string]interface{}, parent *hclwrite.Body, parentBlock string) {
	keys := make([]string, 0, len(structData))
	for k := range structData {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, block := range keys {
		if _, ok := schemaBlock.NestedBlocks[block]; ok {
			if schemaBlock.NestedBlocks[block].NestingMode == "list" || schemaBlock.NestedBlocks[block].NestingMode == "set" {
				child := hclwrite.NewBlock(block, []string{})
				switch s := structData[block].(type) {
				case []map[string]interface{}:
					for _, nestedItem := range s {
						stepChild := hclwrite.NewBlock(block, []string{})
						processBlocks(schemaBlock.NestedBlocks[block].Block, nestedItem, stepChild.Body(), block)
						if len(stepChild.Body().Attributes()) != 0 || len(stepChild.Body().Blocks()) != 0 {
							parent.AppendBlock(stepChild)
						}
					}
				case map[string]interface{}:
					processBlocks(schemaBlock.NestedBlocks[block].Block, s, child.Body(), block)
				case []interface{}:
					for _, nestedItem := range s {
						stepChild := hclwrite.NewBlock(block, []string{})
						processBlocks(schemaBlock.NestedBlocks[block].Block, nestedItem.(map[string]interface{}), stepChild.Body(), block)
						if len(stepChild.Body().Attributes()) != 0 || len(stepChild.Body().Blocks()) != 0 {
							parent.AppendBlock(stepChild)
						}
					}
				default:
					log.Debugf("unable to generate recursively nested blocks for %T", s)
				}
				if len(child.Body().Attributes()) != 0 || len(child.Body().Blocks()) != 0 {
					parent.AppendBlock(child)
				}
			}
		} else {
			if parentBlock == "" && block == "id" {
				continue
			}
			if _, ok := schemaBlock.Attributes[block]; ok && (schemaBlock.Attributes[block].Optional || schemaBlock.Attributes[block].Required) {
				writeAttrLine(block, structData[block], parentBlock, parent)
			}
		}
	}
}

// writeAttrLine outputs a line of HCL configuration with a configurable depth
// for known types.
func writeAttrLine(key string, value interface{}, parentName string, body *hclwrite.Body) {
	switch values := value.(type) {
	case []map[string]interface{}:
		var childCty []cty.Value
		for _, item := range value.([]map[string]interface{}) {
			mapCty := make(map[string]cty.Value)
			for k, v := range item {
				mapCty[k] = cty.StringVal(v.(string))
			}
			childCty = append(childCty, cty.MapVal(mapCty))
		}
		body.SetAttributeValue(key, cty.ListVal(childCty))
	case map[string]interface{}:

		sortedKeys := make([]string, 0, len(values))
		for k := range values {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)

		ctyMap := make(map[string]cty.Value)
		for _, v := range sortedKeys {
			switch val := values[v].(type) {
			case string:
				ctyMap[v] = cty.StringVal(val)
			case float64:
				ctyMap[v] = cty.NumberFloatVal(val)
			}
		}
		body.SetAttributeValue(key, cty.ObjectVal(ctyMap))
	case []interface{}:
		var stringItems []string
		var intItems []int
		var interfaceItems []map[string]interface{}

		for _, item := range value.([]interface{}) {
			switch item := item.(type) {
			case string:
				stringItems = append(stringItems, item)
			case map[string]interface{}:
				interfaceItems = append(interfaceItems, item)
			case float64:
				intItems = append(intItems, int(item))
			}
		}
		if len(stringItems) > 0 {
			writeAttrLine(key, stringItems, parentName, body)
		}

		if len(intItems) > 0 {
			writeAttrLine(key, intItems, parentName, body)
		}

		if len(interfaceItems) > 0 {
			writeAttrLine(key, interfaceItems, parentName, body)
		}
	case []int:
		var vals []cty.Value
		for _, i := range value.([]int) {
			vals = append(vals, cty.NumberIntVal(int64(i)))
		}
		body.SetAttributeValue(key, cty.ListVal(vals))
	case []string:
		var items []string
		items = append(items, value.([]string)...)
		if len(items) > 0 {
			var vals []cty.Value
			for _, item := range items {
				vals = append(vals, cty.StringVal(item))
			}
			body.SetAttributeValue(key, cty.ListVal(vals))
		} else {
			body.SetAttributeValue(key, cty.ListValEmpty(cty.String))
		}
	case string:
		if parentName == "query" && key == "value" && value == "" {
			body.SetAttributeValue(key, cty.StringVal(""))
		}

		if value != "" {
			body.SetAttributeValue(key, cty.StringVal(value.(string)))
		}
	case int:
		body.SetAttributeValue(key, cty.NumberIntVal(int64(value.(int))))
	case float64:
		body.SetAttributeValue(key, cty.NumberFloatVal(value.(float64)))
	case bool:
		body.SetAttributeValue(key, cty.BoolVal(value.(bool)))
	default:
		log.Debugf("got unknown attribute configuration: key %s, value %v, value type %T", key, value, value)
	}
}

// boolToEnabledOrDisabled outputs a string representation of a boolean in the form of `enabled` or `disabled`.
func boolToEnabledOrDisabled(value bool) string {
	if value {
		return "enabled"
	}
	return "disabled"
}
