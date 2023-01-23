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

func buildBlocks(schemaBlock *tfjson.SchemaBlock, structData map[string]interface{}) string {
	attributes := make(map[string]interface{})
	flatten("", structData, attributes)
	return nestBlocks(schemaBlock, attributes, "")
}

// flattening the API response into a queryable map, that follows the nesting layers of our TF resource schema
// i.e. the following entry: attributes["rules.2.action_parameters.edge_ttl.status_code_ttl.0.status_code_range"]
// contains map[from:100 to:200].
func flatten(keyPrefix string, structData interface{}, attributes map[string]interface{}) {
	switch structData.(type) {
	case map[string]interface{}:
		for k1, e1 := range structData.(map[string]interface{}) {
			switch e1.(type) {
			case map[string]interface{}:
				putAttribute(attributes, keyPrefix, k1, e1)
				if keyPrefix != "" {
					flatten(fmt.Sprintf("%v.%v", keyPrefix, k1), e1, attributes)
				} else {
					flatten(fmt.Sprintf("%v", k1), e1, attributes)
				}
			case []interface{}:
				putAttribute(attributes, keyPrefix, k1, e1)
				for k2, e2 := range e1.([]interface{}) {
					putAttribute(attributes, keyPrefix, fmt.Sprintf("%v.%v", k1, k2), e2)
					if keyPrefix != "" {
						flatten(fmt.Sprintf("%v.%v", keyPrefix, fmt.Sprintf("%v.%v", k1, k2)), e2, attributes)
					} else {
						flatten(fmt.Sprintf("%v.%v", k1, k2), e2, attributes)
					}
				}
			case []map[string]interface{}:
				putAttribute(attributes, keyPrefix, k1, e1)
				for k2, e2 := range e1.([]map[string]interface{}) {
					putAttribute(attributes, keyPrefix, fmt.Sprintf("%v.%v", k1, k2), e2)
					if keyPrefix != "" {
						flatten(fmt.Sprintf("%v.%v", keyPrefix, fmt.Sprintf("%v.%v", k1, k2)), e2, attributes)
					} else {
						flatten(fmt.Sprintf("%v.%v", k1, k2), e2, attributes)
					}
				}
			default:
				putAttribute(attributes, keyPrefix, k1, e1)
			}
		}
	case []interface{}:
		for index, value := range structData.([]interface{}) {
			flatten(fmt.Sprintf("%v.%v", keyPrefix, index), value, attributes)
			if keyPrefix != "" {
				flatten(fmt.Sprintf("%v.%v", keyPrefix, index), value, attributes)
			} else {
				flatten(fmt.Sprintf("%v", index), value, attributes)
			}
		}
	}
}

func putAttribute(attributes map[string]interface{}, keyPrefix, key string, element interface{}) {
	if keyPrefix != "" {
		attributes[fmt.Sprintf("%v.%v", keyPrefix, key)] = element
	} else {
		attributes[fmt.Sprintf("%v", key)] = element
	}
}

func nestBlocks(schemaBlock *tfjson.SchemaBlock, attributes map[string]interface{}, queryPrefix string) string {
	output := ""

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

			// Build a query from block names and array iterators to match the schma level we are at.
			// i.e. rules.2.action_parameters.edge_ttl.status_code_ttl.0.status_code_range.
			var query string
			if queryPrefix == "" {
				query = block
			} else {
				query = fmt.Sprintf("%v.%v", queryPrefix, block)
			}

			// Let's query that attribute from our map.
			attribute := attributes[query]
			if attribute == nil {
				continue
			}

			switch attribute.(type) {
			case []interface{}:
				for i, attr := range attribute.([]interface{}) {
					indexedQuery := fmt.Sprintf("%v.%v", query, i)
					nestedBlockOutput := nestBlocks(schemaBlock.NestedBlocks[block].Block, attributes, indexedQuery)
					repeatedBlockOutput := writeNestedBlock(sortedInnerAttributes, schemaBlock.NestedBlocks[block].Block, attr.(map[string]interface{}))
					appendBlock(&output, block, nestedBlockOutput, repeatedBlockOutput)
				}
			case []map[string]interface{}:
				for i, attr := range attribute.([]map[string]interface{}) {
					indexedQuery := fmt.Sprintf("%v.%v", query, i)
					nestedBlockOutput := nestBlocks(schemaBlock.NestedBlocks[block].Block, attributes, indexedQuery)
					repeatedBlockOutput := writeNestedBlock(sortedInnerAttributes, schemaBlock.NestedBlocks[block].Block, attr)
					appendBlock(&output, block, nestedBlockOutput, repeatedBlockOutput)
				}
			case map[string]interface{}:
				nestedBlockOutput := nestBlocks(schemaBlock.NestedBlocks[block].Block, attributes, query)
				repeatedBlockOutput := writeNestedBlock(sortedInnerAttributes, schemaBlock.NestedBlocks[block].Block, attribute.(map[string]interface{}))
				appendBlock(&output, block, nestedBlockOutput, repeatedBlockOutput)
			default:
				log.Debugf("unexpected attribute struct type %T for block %s", attribute, block)
			}
		} else {
			log.Debugf("nested mode %q for %s not recognised", schemaBlock.NestedBlocks[block].NestingMode, block)
		}
	}

	return output
}

func appendBlock(output *string, block, nestedBlockOutput, repeatedBlockOutput string) {
	if repeatedBlockOutput != "" {
		*output += block + " {\n"
		*output += repeatedBlockOutput

		if nestedBlockOutput != "" {
			*output += nestedBlockOutput
		}
		*output += "}\n"
	} else if repeatedBlockOutput == "" && nestedBlockOutput != "" {
		*output += block + " {\n"

		if nestedBlockOutput != "" {
			*output += nestedBlockOutput
		}
		*output += "}\n"
	}
}

func writeNestedBlock(attributes []string, schemaBlock *tfjson.SchemaBlock, attrStruct map[string]interface{}) string {
	nestedBlockOutput := ""

	for _, attrName := range attributes {
		ty := schemaBlock.Attributes[attrName].AttributeType

		switch {
		case ty.IsPrimitiveType():
			switch ty {
			case cty.String, cty.Bool, cty.Number:
				nestedBlockOutput += writeAttrLine(attrName, attrStruct[attrName], false)
			default:
				log.Debugf("unexpected primitive type %q", ty.FriendlyName())
			}
		case ty.IsListType(), ty.IsSetType():
			nestedBlockOutput += writeAttrLine(attrName, attrStruct[attrName], true)
		case ty.IsMapType():
			nestedBlockOutput += writeAttrLine(attrName, attrStruct[attrName], false)
		default:
			log.Debugf("unexpected nested type %T for %s", ty, attrName)
		}
	}

	return nestedBlockOutput
}

// writeAttrLine outputs a line of HCL configuration with a configurable depth
// for known types.
func writeAttrLine(key string, value interface{}, usedInBlock bool) string {
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
				s += writeAttrLine(fmt.Sprintf("\"%s\"", v), values[v], false)
			} else {
				s += writeAttrLine(v, values[v], false)
			}
		}

		if usedInBlock {
			if s != "" {
				return fmt.Sprintf("%s {\n%s}\n", key, s)
			}
		} else {
			if s != "" {
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
			return writeAttrLine(key, stringItems, false)
		}

		if len(intItems) > 0 {
			return writeAttrLine(key, intItems, false)
		}

		if len(interfaceItems) > 0 {
			return writeAttrLine(key, interfaceItems, false)
		}

	case []map[string]interface{}:
		var stringyInterfaces []string
		var op string
		var mapLen = len(value.([]map[string]interface{}))
		for i, item := range value.([]map[string]interface{}) {
			// Use an empty key to prevent rendering the key
			op = writeAttrLine("", item, true)
			// if condition handles adding new line for just the last element
			if i != mapLen-1 {
				op = strings.TrimRight(op, "\n")
			}
			stringyInterfaces = append(stringyInterfaces, op)
		}
		return fmt.Sprintf("%s = [ \n%s ]\n", key, strings.Join(stringyInterfaces, ",\n"))

	case []int:
		stringyInts := []string{}
		for _, int := range value.([]int) {
			stringyInts = append(stringyInts, fmt.Sprintf("%d", int))
		}
		return fmt.Sprintf("%s = [ %s ]\n", key, strings.Join(stringyInts, ", "))
	case []string:
		var items []string
		for _, item := range value.([]string) {
			items = append(items, fmt.Sprintf("%q", item))
		}
		if len(items) > 0 {
			return fmt.Sprintf("%s = [ %s ]\n", key, strings.Join(items, ", "))
		}
	case string:
		if value != "" {
			return fmt.Sprintf("%s = %q\n", key, value)
		}
	case int:
		return fmt.Sprintf("%s = %d\n", key, value)
	case float64:
		return fmt.Sprintf("%s = %0.f\n", key, value)
	case bool:
		return fmt.Sprintf("%s = %t\n", key, value)
	default:
		log.Debugf("got unknown attribute configuration: key %s, value %v, value type %T", key, value, value)
		return ""
	}
	return ""
}
