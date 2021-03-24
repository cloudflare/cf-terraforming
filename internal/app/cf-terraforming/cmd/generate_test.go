package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_writeAttrLine(t *testing.T) {
	var listOfString []interface{}
	listOfString = append(listOfString, "b", "c", "d")

	tests := map[string]struct {
		key   string
		value interface{}
		depth int
		want  string
	}{
		"value is string":          {key: "a", value: "b", depth: 0, want: fmt.Sprintf("a = %q\n", "b")},
		"value is int":             {key: "a", value: 1, depth: 0, want: "a = 1\n"},
		"value is float":           {key: "a", value: 1.0, depth: 0, want: "a = 1\n"},
		"value is bool":            {key: "a", value: true, depth: 0, want: "a = true\n"},
		"value is list of strings": {key: "a", value: listOfString, depth: 0, want: "a = [ \"b\", \"c\", \"d\" ]\n"},
		"value is nil":             {key: "a", value: nil, depth: 0, want: ""},

		"depth is 0": {key: "a", value: "b", depth: 0, want: fmt.Sprintf("a = %q\n", "b")},
		"depth is 6": {key: "a", value: "b", depth: 6, want: fmt.Sprintf("      a = %q\n", "b")},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := writeAttrLine(tc.key, tc.value, tc.depth)
			assert.Equal(t, got, tc.want)
		})
	}
}
