package cmd

import (
	"strings"
	"text/template"
)

func replace(input, from, to string) string {
	return strings.Replace(input, from, to, -1)
}

var templateFuncMap = template.FuncMap{
	"replace": replace,
}
