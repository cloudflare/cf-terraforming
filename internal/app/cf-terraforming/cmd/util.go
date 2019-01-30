package cmd

import (
	"fmt"
	"reflect"
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

func isMap(i interface{}) bool {
	return (reflect.ValueOf(i).Kind() == reflect.Map)
}

func quoteIfString(i interface{}) interface{} {
	if reflect.ValueOf(i).Kind() == reflect.String {
		return fmt.Sprintf("\"%v\"", i)
	} else {
		return i
	}
}

var templateFuncMap = template.FuncMap{
	"replace":       replace,
	"isMap":         isMap,
	"quoteIfString": quoteIfString,
}
