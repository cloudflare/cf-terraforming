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
		log.Debug("--account and --zone are mutually exclusive, support for both is deprecated")
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

// nestBlocks takes a schema and generates all of the appropriate nesting of any
// top-level blocks as well as nested lists or sets.
func nestBlocks(schemaBlock *tfjson.SchemaBlock, structData map[string]interface{}, parent *hclwrite.Body) {

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

			for attrName, attrConfig := range schemaBlock.NestedBlocks[block].Block.Attributes {
				if attrConfig.Computed && !attrConfig.Optional {
					schemaBlock.NestedBlocks[block].Block.Attributes[attrName].AttributeType = cty.NilType
				}
			}

			var currentNode *hclwrite.Body
			// If the attribute we're looking at has further nesting, we'll
			// recursively call nestBlocks.
			if block == "status_code_ttl" {
				fmt.Println("DEBUG")
			}
			if len(schemaBlock.NestedBlocks[block].Block.NestedBlocks) > 0 {
				if s, ok := structData[block]; ok {
					switch s.(type) {
					case map[string]interface{}:
						currentNode = parent.AppendNewBlock(block, []string{}).Body()
						nestBlocks(schemaBlock.NestedBlocks[block].Block, s.(map[string]interface{}), currentNode)
					case []interface{}:
						for _, nestedItem := range s.([]interface{}) {
							child := parent.AppendNewBlock(block, []string{}).Body()
							nestBlocks(schemaBlock.NestedBlocks[block].Block, nestedItem.(map[string]interface{}), child)
							processAttributes(schemaBlock, block, sortedInnerAttributes, structData, child)
						}
						return
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
					if currentNode == nil {
						currentNode = parent.AppendNewBlock(block, []string{}).Body()
					}

					writeNestedBlock(sortedInnerAttributes, schemaBlock.NestedBlocks[block].Block, attrStruct, currentNode)
				}

			// Case for if the inner block's attributes are a list of map interfaces,
			// in which case this should be treated as a duplicating block.
			case []map[string]interface{}:
				for _, v := range attrStruct {
					if attrStruct != nil {

						writeNestedBlock(sortedInnerAttributes, schemaBlock.NestedBlocks[block].Block, v, parent.AppendNewBlock(block, []string{}).Body())
					}
				}

			// Case for duplicated blocks that commonly end up as an array or list at
			// the API level.
			case []interface{}:
				for _, v := range attrStruct {
					if attrStruct != nil {
						writeNestedBlock(sortedInnerAttributes, schemaBlock.NestedBlocks[block].Block, v.(map[string]interface{}), parent.AppendNewBlock(block, []string{}).Body())
					}
				}

			default:
				log.Debugf("unexpected attribute struct type %T for block %s", attrStruct, block)
			}
		} else {
			log.Debugf("nested mode %q for %s not recognised", schemaBlock.NestedBlocks[block].NestingMode, block)
		}
	}
}

func processAttributes(schemaBlock *tfjson.SchemaBlock, block string, sortedInnerAttributes []string, structData map[string]interface{}, body *hclwrite.Body) {
	switch attrStruct := structData[block].(type) {
	// Case for if the inner block's attributes are a map of interfaces, in
	// which case we can directly add them to the config.
	case map[string]interface{}:
		if attrStruct != nil {
			writeNestedBlock(sortedInnerAttributes, schemaBlock.NestedBlocks[block].Block, attrStruct, body)
		}

	// Case for if the inner block's attributes are a list of map interfaces,
	// in which case this should be treated as a duplicating block.
	case []map[string]interface{}:
		for _, v := range attrStruct {
			if attrStruct != nil {
				writeNestedBlock(sortedInnerAttributes, schemaBlock.NestedBlocks[block].Block, v, body)
			}
		}

	// Case for duplicated blocks that commonly end up as an array or list at
	// the API level.
	case []interface{}:
		for _, v := range attrStruct {
			if attrStruct != nil {
				bodyBlocks := body.Blocks()
				for _, b := range bodyBlocks {
					name := b.Type()
					attrMap := make(map[string]interface{}, 0)
					for k, val := range b.Body().Attributes() {
						token := val.Expr().BuildTokens(nil)
						value := string(token.Bytes())
						attrMap[k] = value
					}
					if fmt.Sprint(v.(map[string]interface{})[name]) == fmt.Sprint(attrMap) {
						writeNestedBlock(sortedInnerAttributes, schemaBlock.NestedBlocks[block].Block, v.(map[string]interface{}), body)
					}
				}
			}
		}

	default:
		log.Debugf("unexpected attribute struct type %T for block %s", attrStruct, block)
	}
}

func writeNestedBlock(attributes []string, schemaBlock *tfjson.SchemaBlock, attrStruct map[string]interface{}, body *hclwrite.Body) string {
	nestedBlockOutput := ""

	for _, attrName := range attributes {

		ty := schemaBlock.Attributes[attrName].AttributeType

		switch {
		case ty.IsPrimitiveType():
			switch ty {
			case cty.String, cty.Bool, cty.Number:
				nestedBlockOutput += writeAttrLine(attrName, attrStruct[attrName], false, body)
			default:
				log.Debugf("unexpected primitive type %q", ty.FriendlyName())
			}
		case ty.IsListType(), ty.IsSetType():
			nestedBlockOutput += writeAttrLine(attrName, attrStruct[attrName], true, body)
		case ty.IsMapType():
			nestedBlockOutput += writeAttrLine(attrName, attrStruct[attrName], false, body)
		default:
			log.Debugf("unexpected nested type %T for %s", ty, attrName)
		}
	}

	return nestedBlockOutput
}

// writeAttrLine outputs a line of HCL configuration with a configurable depth
// for known types.
func writeAttrLine(key string, value interface{}, usedInBlock bool, body *hclwrite.Body) string {
	switch values := value.(type) {
	case map[string]interface{}:
		sortedKeys := make([]string, 0, len(values))
		for k := range values {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)

		s := ""
		for _, v := range sortedKeys {
			// check if our key has an integer in the string. If it does we need to wrap it with quotes.
			if hasNumber(v) {
				s += writeAttrLine(fmt.Sprintf("\"%s\"", v), values[v], false, body)
			} else {
				s += writeAttrLine(v, values[v], false, body)
			}
		}

		if usedInBlock {
			if s != "" {
				body.AppendNewBlock(key, []string{}).Body().SetAttributeValue(s, cty.NilVal)
				return fmt.Sprintf("%s {\n%s}\n", key, s)
			}
		} else {
			if s != "" {
				body.SetAttributeValue(key, cty.StringVal(s))
				return fmt.Sprintf("%s = {\n%s}\n", key, s)
			}
		}
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
			return writeAttrLine(key, stringItems, false, body)
		}

		if len(intItems) > 0 {
			return writeAttrLine(key, intItems, false, body)
		}

		if len(interfaceItems) > 0 {
			return writeAttrLine(key, interfaceItems, false, body)
		}

	case []map[string]interface{}:
		var stringyInterfaces []string
		var op string
		var mapLen = len(value.([]map[string]interface{}))
		for i, item := range value.([]map[string]interface{}) {
			// Use an empty key to prevent rendering the key
			op = writeAttrLine("", item, true, body)
			// if condition handles adding new line for just the last element
			if i != mapLen-1 {
				op = strings.TrimRight(op, "\n")
			}
			stringyInterfaces = append(stringyInterfaces, op)
		}
		body.SetAttributeValue(key, cty.StringVal(fmt.Sprintf("[%s]", strings.Join(stringyInterfaces, ",\n"))))
		return fmt.Sprintf("%s = [ \n%s ]\n", key, strings.Join(stringyInterfaces, ",\n"))

	case []int:
		stringyInts := []string{}
		for _, int := range value.([]int) {
			stringyInts = append(stringyInts, fmt.Sprintf("%d", int))
		}
		body.SetAttributeValue(key, cty.StringVal(fmt.Sprintf("[%s]", strings.Join(stringyInts, ", "))))
		return fmt.Sprintf("%s = [ %s ]\n", key, strings.Join(stringyInts, ", "))
	case []string:
		var items []string
		for _, item := range value.([]string) {
			items = append(items, fmt.Sprintf("%q", item))
		}
		if len(items) > 0 {
			body.SetAttributeValue(key, cty.StringVal(fmt.Sprintf("[%s]", strings.Join(items, ", "))))
			return fmt.Sprintf("%s = [ %s ]\n", key, strings.Join(items, ", "))
		}
	case string:
		if value != "" {
			body.SetAttributeValue(key, cty.StringVal(value.(string)))
			return fmt.Sprintf("%s = %q\n", key, value)
		}
	case int:
		body.SetAttributeValue(key, cty.NumberIntVal(int64(value.(int))))
		return fmt.Sprintf("%s = %d\n", key, value)
	case float64:
		body.SetAttributeValue(key, cty.NumberFloatVal(value.(float64)))
		return fmt.Sprintf("%s = %0.f\n", key, value)
	case bool:
		body.SetAttributeValue(key, cty.BoolVal(value.(bool)))
		return fmt.Sprintf("%s = %t\n", key, value)
	default:
		log.Debugf("got unknown attribute configuration: key %s, value %v, value type %T", key, value, value)
		return ""
	}
	return ""
}
