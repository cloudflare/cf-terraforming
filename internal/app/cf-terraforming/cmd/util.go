package cmd

import (
	"strings"
	"text/template"
)

func replace(input, from, to string) string {
	return strings.Replace(input, from, to, -1)
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

var templateFuncMap = template.FuncMap{
	"replace": replace,
}
