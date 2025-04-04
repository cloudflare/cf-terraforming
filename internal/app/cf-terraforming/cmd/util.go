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

	cfv0 "github.com/cloudflare/cloudflare-go"
	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/hashicorp/hcl/v2/hclwrite"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
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
func testDataFile(version, filename string) string {
	filename = strings.TrimSuffix(filename, "/")

	dirname, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dir, err := os.Open(filepath.Join(dirname, fmt.Sprintf("../../../../testdata/terraform/%s", version)))
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

	var options []cfv0.Option

	if hostname != "" {
		options = append(options, cfv0.BaseURL("https://"+hostname+"/client/v4"))
	}

	if verbose {
		options = append(options, cfv0.Debug(true))
	}

	var err error

	// Don't initialise a client in CI as this messes with VCR and the ability to
	// mock out the HTTP interactions.
	if os.Getenv("CI") != "true" {
		var useToken = apiToken != ""

		if useToken {
			apiV0, err = cfv0.NewWithAPIToken(apiToken, options...)
			api = cloudflare.NewClient(option.WithAPIToken(apiToken))
		} else {
			apiV0, err = cfv0.New(apiKey, apiEmail, options...)
			api = cloudflare.NewClient(option.WithAPIKey(apiKey), option.WithAPIEmail(apiEmail))
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

// flattenAttrMap takes a list of attributes defined as a list of maps comprising {"id": "attrId", "value": "attrValue"}
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
	if body == nil || value == nil {
		log.Debug("body or value is nil")
		return
	}

	switch values := value.(type) {
	case []map[string]interface{}:
		// Use tuple approach for heterogeneous maps
		var tupleValues []cty.Value
		for _, item := range values {
			mapCty := make(map[string]cty.Value)
			for k, v := range item {
				mapCty[k] = processExpression(v)
			}
			tupleValues = append(tupleValues, cty.ObjectVal(mapCty))
		}
		body.SetAttributeValue(key, cty.TupleVal(tupleValues))
	case map[string]interface{}:
		ctyMap := make(map[string]cty.Value)

		// Sort keys for consistent output
		sortedKeys := make([]string, 0, len(values))
		for k := range values {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)

		for _, k := range sortedKeys {
			ctyMap[k] = processExpression(values[k])
		}
		body.SetAttributeValue(key, cty.ObjectVal(ctyMap))
	case []interface{}:
		if len(values) == 0 {
			body.SetAttributeValue(key, cty.EmptyTupleVal)
			return
		}

		// Convert all slice elements using processExpression for consistency
		var tupleValues []cty.Value
		for _, item := range values {
			tupleValues = append(tupleValues, processExpression(item))
		}
		body.SetAttributeValue(key, cty.TupleVal(tupleValues))
	case []int:
		var vals []cty.Value
		for _, i := range values {
			vals = append(vals, cty.NumberIntVal(int64(i)))
		}
		body.SetAttributeValue(key, cty.TupleVal(vals))
	case []string:
		if len(values) > 0 {
			var vals []cty.Value
			for _, item := range values {
				vals = append(vals, cty.StringVal(item))
			}
			body.SetAttributeValue(key, cty.TupleVal(vals))
		} else {
			body.SetAttributeValue(key, cty.EmptyTupleVal)
		}
	case string:
		if parentName == "query" && key == "value" && value == "" {
			body.SetAttributeValue(key, cty.StringVal(""))
		}
		if value != "" {
			body.SetAttributeValue(key, cty.StringVal(values))
		}
	case int:
		body.SetAttributeValue(key, cty.NumberIntVal(int64(values)))
	case float64:
		body.SetAttributeValue(key, cty.NumberFloatVal(values))
	case bool:
		body.SetAttributeValue(key, cty.BoolVal(values))
	default:
		fmt.Printf("Warning: Unknown attribute type: key %s, value %v, value type %T\n", key, value, value)
		// Convert unknown types to string representation
		body.SetAttributeValue(key, cty.StringVal(fmt.Sprintf("%v", value)))
	}
}

// Process any expression into its appropriate cty.Value and also modified to use TupleVal consistently.
func processExpression(val interface{}) cty.Value {
	if val == nil {
		return cty.NullVal(cty.DynamicPseudoType)
	}

	switch v := val.(type) {
	case string:
		return cty.StringVal(v)
	case int:
		return cty.NumberIntVal(int64(v))
	case float64:
		return cty.NumberFloatVal(v)
	case bool:
		return cty.BoolVal(v)
	case []string:
		var vals []cty.Value
		for _, s := range v {
			vals = append(vals, cty.StringVal(s))
		}
		return cty.TupleVal(vals)
	case []int:
		var vals []cty.Value
		for _, i := range v {
			vals = append(vals, cty.NumberIntVal(int64(i)))
		}
		return cty.TupleVal(vals)
	case []interface{}:
		if len(v) == 0 {
			return cty.EmptyTupleVal
		}

		var vals []cty.Value
		for _, item := range v {
			vals = append(vals, processExpression(item))
		}
		return cty.TupleVal(vals)
	case map[string]interface{}:
		ctyMap := make(map[string]cty.Value)
		// Sort keys for consistent output
		sortedKeys := make([]string, 0, len(v))
		for k := range v {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)

		for _, k := range sortedKeys {
			ctyMap[k] = processExpression(v[k])
		}
		return cty.ObjectVal(ctyMap)
	case []map[string]interface{}:
		var vals []cty.Value
		for _, m := range v {
			// Convert map to object
			objMap := make(map[string]cty.Value)
			for mk, mv := range m {
				objMap[mk] = processExpression(mv)
			}
			vals = append(vals, cty.ObjectVal(objMap))
		}
		return cty.TupleVal(vals)
	default:
		return cty.StringVal(fmt.Sprintf("%v", val))
	}
}

// boolToEnabledOrDisabled outputs a string representation of a boolean in the form of `enabled` or `disabled`.
func boolToEnabledOrDisabled(value bool) string {
	if value {
		return "enabled"
	}
	return "disabled"
}

// transformToCollection takes a JSON payload that is a singlular resource but
// operates as a `list` endpoint and transforms it into a JSON to correctly
// handle the output.
func transformToCollection(value string) string {
	return fmt.Sprintf("[%s]", value)
}

// modifyResponsePayload takes the current resource and the `gjson.Result`
// to run arbitrary modifications to the JSON before passing it to be overlayed
// the provider schema.
func modifyResponsePayload(resourceName string, value gjson.Result) string {
	output := value.String()

	switch resourceName {
	case "cloudflare_zero_trust_organization":
		output = transformToCollection(output)
	}

	return output
}
