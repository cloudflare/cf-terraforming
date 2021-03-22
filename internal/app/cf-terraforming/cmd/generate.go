package cmd

import (
	"github.com/spf13/cobra"

	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
	"github.com/zclconf/go-cty/cty"
)

var (
	resourceType string

	// eventually, this will come from the API
	jsonPayload = []byte(`
	{
		"id": "372e67954025e0ba6aaa6d586b9e0b59",
		"type": "A",
		"name": "example.com",
		"content": "198.51.100.4",
		"proxiable": true,
		"proxied": false,
		"ttl": 120,
		"locked": false,
		"zone_id": "023e105f4ecef8ad9ca31a8372d0c353",
		"zone_name": "example.com",
		"created_on": "2014-01-01T05:20:00.12345Z",
		"modified_on": "2014-01-01T05:20:00.12345Z",
		"data": {},
		"meta": {
			"auto_added": true,
			"source": "primary"
		}
	}`)
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

		var record cloudflare.DNSRecord
		json.Unmarshal(jsonPayload, &record)

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
		tf, err := tfexec.NewTerraform(workingDir, execPath)
		if err != nil {
			log.Fatal(err)
		}

		err = tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.LockTimeout("60s"))
		if err != nil {
			log.Fatal(err)
		}

		ps, err := tf.ProvidersSchema(context.Background())
		s := ps.Schemas["registry.terraform.io/cloudflare/cloudflare"]
		if s == nil {
			log.Fatal("failed to detect provider installation")
		}

		r := s.ResourceSchemas[*&resourceType]

		output := ""
		output += fmt.Sprintf(`resource "%s" "some_generated_name" {`+"\n", *&resourceType)

		for attrName, attrConfig := range r.Block.Attributes {
			t := reflect.TypeOf(record)

			// Iterate over all available fields and read the tag value
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				tag := field.Tag.Get("json")
				jsonTag := strings.Split(tag, ",")[0]

				if jsonTag == attrName {
					r := reflect.ValueOf(record)
					f := reflect.Indirect(r).FieldByName(field.Name)

					ty := attrConfig.AttributeType
					switch {
					case ty.IsPrimitiveType():
						switch ty {
						case cty.String:
							if attrName == "created_on" || attrName == "modified_on" {
								t := f.Interface().(time.Time)
								formatedTime := t.Format(time.RFC3339Nano)
								output += writeAttrLine(attrName, formatedTime, 2)
							} else {
								output += writeAttrLine(attrName, f.Interface().(string), 2)
							}
						case cty.Bool:
							output += writeAttrLine(jsonTag, f.Interface().(bool), 2)
						case cty.Number:
							output += writeAttrLine(attrName, f.Interface().(int), 2)
						default:
							log.Warnf("unexpected primitive type %q", ty.FriendlyName())
						}
					case ty.IsCollectionType():
						switch {
						case ty.IsListType(), ty.IsSetType():
							if f.Interface() == nil || len(f.Interface().([]string)) == 0 {
								continue
							}

							items := f.Interface().([]string)
							for k, v := range items {
								items[k] = fmt.Sprintf("%q", v)
							}
							output += fmt.Sprintf("  %s = [ ", attrName)
							output += strings.Join(items, ", ")
							output += fmt.Sprintf(" ]\n")
						case ty.IsMapType():
							if f.Interface() == nil {
								continue
							}

							// output += buildNestedMap(attrName, f, r)

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
			}
		}

		output += "}\n"
		fmt.Println(output)
	},
}
// writeAttrLine outputs a line of HCL configuration with a configurable depth
// for known types.
func writeAttrLine(key string, value interface{}, depth int) string {
	switch reflect.TypeOf(value).String() {
	case "string":
		return fmt.Sprintf("%s%s = %q\n", strings.Repeat(" ", depth), key, value)
	case "int":
		return fmt.Sprintf("%s%s = %d\n", strings.Repeat(" ", depth), key, value)
	case "float64":
		return fmt.Sprintf("%s%s = %0.f\n", strings.Repeat(" ", depth), key, value)
	case "bool":
		return fmt.Sprintf("%s%s = %t\n", strings.Repeat(" ", depth), key, value)
	default:
		log.Debugf("got unknown attribute configuration: key %s, value %v, value type %T", key, value, value)
		return "n/a"
	}
}
