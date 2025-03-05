package cmd

import (
	"sort"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestProcessExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected cty.Value
	}{
		{
			name:     "nil value",
			input:    nil,
			expected: cty.NullVal(cty.DynamicPseudoType),
		},
		{
			name:     "string value",
			input:    "test",
			expected: cty.StringVal("test"),
		},
		{
			name:     "int value",
			input:    42,
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "float value",
			input:    3.14,
			expected: cty.NumberFloatVal(3.14),
		},
		{
			name:     "bool value",
			input:    true,
			expected: cty.BoolVal(true),
		},
		{
			name:  "string slice",
			input: []string{"a", "b", "c"},
			expected: cty.TupleVal([]cty.Value{
				cty.StringVal("a"),
				cty.StringVal("b"),
				cty.StringVal("c"),
			}),
		},
		{
			name:  "int slice",
			input: []int{1, 2, 3},
			expected: cty.TupleVal([]cty.Value{
				cty.NumberIntVal(1),
				cty.NumberIntVal(2),
				cty.NumberIntVal(3),
			}),
		},
		{
			name:     "empty interface slice",
			input:    []interface{}{},
			expected: cty.EmptyTupleVal,
		},
		{
			name:  "mixed interface slice",
			input: []interface{}{"a", 1, true},
			expected: cty.TupleVal([]cty.Value{
				cty.StringVal("a"),
				cty.NumberIntVal(1),
				cty.BoolVal(true),
			}),
		},
		{
			name:  "simple map",
			input: map[string]interface{}{"key": "value"},
			expected: cty.ObjectVal(map[string]cty.Value{
				"key": cty.StringVal("value"),
			}),
		},
		{
			name: "complex map",
			input: map[string]interface{}{
				"str":  "value",
				"num":  42,
				"bool": true,
				"list": []string{"a", "b"},
			},
			expected: cty.ObjectVal(map[string]cty.Value{
				"str":  cty.StringVal("value"),
				"num":  cty.NumberIntVal(42),
				"bool": cty.BoolVal(true),
				"list": cty.TupleVal([]cty.Value{
					cty.StringVal("a"),
					cty.StringVal("b"),
				}),
			}),
		},
		{
			name: "slice of maps",
			input: []map[string]interface{}{
				{"name": "item1", "value": 1},
				{"name": "item2", "value": 2},
			},
			expected: cty.TupleVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"name":  cty.StringVal("item1"),
					"value": cty.NumberIntVal(1),
				}),
				cty.ObjectVal(map[string]cty.Value{
					"name":  cty.StringVal("item2"),
					"value": cty.NumberIntVal(2),
				}),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processExpression(tt.input)
			// For maps, we need to check that all keys are present with expected values
			// rather than exact equality since key order might differ
			if _, ok := tt.input.(map[string]interface{}); ok {
				if !ctyValuesEquivalent(result, tt.expected) {
					t.Errorf("processExpression() = %v, want %v", result, tt.expected)
				}
			} else {
				if !result.RawEquals(tt.expected) {
					t.Errorf("processExpression() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestWriteAttrLine(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		value      interface{}
		parentName string
		expected   string
	}{
		{
			name:       "string value",
			key:        "name",
			value:      "test",
			parentName: "",
			expected:   `name = "test"`,
		},
		{
			name:       "empty string in query",
			key:        "value",
			value:      "",
			parentName: "query",
			expected:   `value = ""`,
		},
		{
			name:       "int value",
			key:        "count",
			value:      42,
			parentName: "",
			expected:   `count = 42`,
		},
		{
			name:       "float value",
			key:        "ratio",
			value:      3.14,
			parentName: "",
			expected:   `ratio = 3.14`,
		},
		{
			name:       "bool value",
			key:        "enabled",
			value:      true,
			parentName: "",
			expected:   `enabled = true`,
		},
		{
			name:       "string slice",
			key:        "tags",
			value:      []string{"a", "b", "c"},
			parentName: "",
			expected:   `tags = ["a", "b", "c"]`,
		},
		{
			name:       "int slice",
			key:        "ports",
			value:      []int{80, 443, 8080},
			parentName: "",
			expected:   `ports = [80, 443, 8080]`,
		},
		{
			name:       "empty string slice",
			key:        "empty_tags",
			value:      []string{},
			parentName: "",
			expected:   `empty_tags = []`,
		},
		{
			name:       "simple map",
			key:        "metadata",
			value:      map[string]interface{}{"app": "service", "env": "prod"},
			parentName: "",
			expected:   `metadata = {app = "service", env = "prod"}`,
		},
		{
			name: "complex map",
			key:  "config",
			value: map[string]interface{}{
				"name":    "app",
				"version": 1,
				"enabled": true,
				"tags":    []string{"web", "api"},
			},
			parentName: "",
			expected:   `config = {enabled = true, name = "app", tags = ["web", "api"], version = 1}`,
		},
		{
			name: "slice of maps",
			key:  "resources",
			value: []map[string]interface{}{
				{"name": "res1", "count": 1},
				{"name": "res2", "count": 2},
			},
			parentName: "",
			expected:   `resources = [{count = 1, name = "res1"}, {count = 2, name = "res2"}]`,
		},
		{
			name:       "empty slice",
			key:        "empty_list",
			value:      []interface{}{},
			parentName: "",
			expected:   `empty_list = []`,
		},
		{
			name:       "mixed slice",
			key:        "mixed",
			value:      []interface{}{"string", 42, true},
			parentName: "",
			expected:   `mixed = ["string", 42, true]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := hclwrite.NewEmptyFile()
			rootBody := f.Body()

			writeAttrLine(tt.key, tt.value, tt.parentName, rootBody)

			result := string(f.Bytes())
			// Trim trailing newline for comparison
			if len(result) > 0 && result[len(result)-1] == '\n' {
				result = result[:len(result)-1]
			}

			// Parse the expected and actual HCL to compare normalized representations
			expectedNormalized := normalizeHCL(t, tt.expected)
			actualNormalized := normalizeHCL(t, result)

			assert.Equal(t, expectedNormalized, actualNormalized)
		})
	}
}

func TestWriteAttrLine_NilCases(t *testing.T) {
	t.Run("nil body", func(t *testing.T) {
		// Should not panic
		writeAttrLine("key", "value", "", nil)
	})

	t.Run("nil value", func(t *testing.T) {
		f := hclwrite.NewEmptyFile()
		rootBody := f.Body()

		// Should not write anything
		writeAttrLine("key", nil, "", rootBody)

		result := string(f.Bytes())
		assert.Equal(t, "", result)
	})
}

// Helper function to normalize HCL by parsing and generating a new HCL file.
func normalizeHCL(t *testing.T, hclString string) string {
	// Parse the HCL content
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL([]byte(hclString), "test.hcl")
	if diags.HasErrors() {
		t.Fatalf("Failed to parse HCL: %s", diags.Error())
	}

	// Create a new HCL file with the same content
	f := hclwrite.NewEmptyFile()
	ctx := &hcl.EvalContext{}

	// Extract attributes and recreate them in a consistent order
	attrs, _ := file.Body.JustAttributes()
	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		attr := attrs[k]
		val, diags := attr.Expr.Value(ctx)
		if diags.HasErrors() {
			t.Fatalf("Failed to evaluate expression: %s", diags.Error())
		}
		f.Body().SetAttributeValue(k, val)
	}

	return string(f.Bytes())
}

// Helper function to compare cty.Value objects, handling maps with potentially different key order.
func ctyValuesEquivalent(a, b cty.Value) bool {
	if a.Type().IsObjectType() && b.Type().IsObjectType() {
		aMap := a.AsValueMap()
		bMap := b.AsValueMap()

		if len(aMap) != len(bMap) {
			return false
		}

		for k, aVal := range aMap {
			bVal, ok := bMap[k]
			if !ok || !ctyValuesEquivalent(aVal, bVal) {
				return false
			}
		}
		return true
	}

	if a.Type().IsTupleType() && b.Type().IsTupleType() {
		aSlice := a.AsValueSlice()
		bSlice := b.AsValueSlice()

		if len(aSlice) != len(bSlice) {
			return false
		}

		for i := range aSlice {
			if !ctyValuesEquivalent(aSlice[i], bSlice[i]) {
				return false
			}
		}
		return true
	}

	return a.RawEquals(b)
}
