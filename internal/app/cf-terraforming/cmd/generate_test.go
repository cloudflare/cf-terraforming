package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

var (
	// listOfString is an example representation of a key where the value is a
	// list of string values.
	//
	// resource "example" "example" {
	//   attr = [ "b", "c", "d"]
	// }
	listOfString = []interface{}{"b", "c", "d"}

	// configBlockOfStrings is an example of where a key is a "block" assignment
	// in HCL.
	//
	// resource "example" "example" {
	//   attr = {
	//     c = "d"
	//     e = "f"
	//   }
	// }
	configBlockOfStrings = map[string]interface{}{
		"c": "d",
		"e": "f",
	}
)

func TestGenerate_writeAttrLine(t *testing.T) {
	tests := map[string]struct {
		key   string
		value interface{}
		depth int
		want  string
	}{
		"value is string":           {key: "a", value: "b", depth: 0, want: fmt.Sprintf("a = %q\n", "b")},
		"value is int":              {key: "a", value: 1, depth: 0, want: "a = 1\n"},
		"value is float":            {key: "a", value: 1.0, depth: 0, want: "a = 1\n"},
		"value is bool":             {key: "a", value: true, depth: 0, want: "a = true\n"},
		"value is list of strings":  {key: "a", value: listOfString, depth: 0, want: "a = [ \"b\", \"c\", \"d\" ]\n"},
		"value is block of strings": {key: "a", value: configBlockOfStrings, depth: 0, want: "a = {\n  c = \"d\"\n  e = \"f\"\n}\n"},
		"value is nil":              {key: "a", value: nil, depth: 0, want: ""},

		"depth is 0": {key: "a", value: "b", depth: 0, want: fmt.Sprintf("a = %q\n", "b")},
		"depth is 6": {key: "a", value: "b", depth: 6, want: fmt.Sprintf("      a = %q\n", "b")},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := writeAttrLine(tc.key, tc.value, tc.depth, false)
			assert.Equal(t, got, tc.want)
		})
	}
}

func TestGenerate_ResourceNotSupported(t *testing.T) {
	_, output, err := executeCommandC(GenerateCmd(), "--resource-type", "notreal")

	if assert.Nil(t, err) {
		assert.Contains(t, output, "\"notreal\" is not yet supported for automatic generation")
	}
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}
