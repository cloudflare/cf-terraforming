package cmd

import (
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// postProcess allows you to perform additional actions on the generated hcl.
func postProcess(f *hclwrite.File, resourceType string) {
	switch resourceType {
	case "cloudflare_stream_live_input", "cloudflare_stream":
		addJSONEncode(f, "meta")
	case "cloudflare_observatory_scheduled_test":
		addURLEncode(f, "url")
	}
}

// addJSONEncode wraps a hcl block with the jsonencode function.
func addJSONEncode(f *hclwrite.File, attributeName string) {
	for _, block := range f.Body().Blocks() {
		if block.Type() != "resource" {
			continue
		}
		if len(block.Labels()) < 1 {
			continue
		}
		if block.Labels()[0] != resourceType {
			continue
		}
		body := block.Body()
		attr := body.GetAttribute(attributeName)
		if attr == nil {
			continue
		}
		exprTokens := attr.Expr().BuildTokens(nil)
		exprText := string(exprTokens.Bytes())

		trimmed := strings.TrimSpace(exprText)
		// Wrap the attribute with jsonencode
		if len(trimmed) > 0 && trimmed[0] == '{' {
			body.RemoveAttribute(attributeName)
			newTokens := hclwrite.Tokens{}
			fnStart := &hclwrite.Token{
				Type:  hclsyntax.TokenIdent,
				Bytes: []byte("jsonencode("),
			}
			newTokens = append(newTokens, fnStart)
			newTokens = append(newTokens, exprTokens...)
			fnEnd := &hclwrite.Token{
				Type:  hclsyntax.TokenCParen,
				Bytes: []byte(")"),
			}
			newTokens = append(newTokens, fnEnd)
			body.SetAttributeRaw(attributeName, newTokens)
		}
	}
}

// addURLEncode wraps a hcl block with the urlencode function.
func addURLEncode(f *hclwrite.File, attributeName string) {
	for _, block := range f.Body().Blocks() {
		if block.Type() != "resource" {
			continue
		}
		if len(block.Labels()) < 1 {
			continue
		}
		if block.Labels()[0] != resourceType {
			continue
		}
		body := block.Body()
		attr := body.GetAttribute(attributeName)
		if attr == nil {
			continue
		}
		exprTokens := attr.Expr().BuildTokens(nil)
		exprText := string(exprTokens.Bytes())

		trimmed := strings.TrimSpace(exprText)
		// Wrap the attribute with jsonencode
		if len(trimmed) > 0 {
			body.RemoveAttribute(attributeName)
			newTokens := hclwrite.Tokens{}
			fnStart := &hclwrite.Token{
				Type:  hclsyntax.TokenIdent,
				Bytes: []byte("urlencode("),
			}
			newTokens = append(newTokens, fnStart)
			newTokens = append(newTokens, exprTokens...)
			fnEnd := &hclwrite.Token{
				Type:  hclsyntax.TokenCParen,
				Bytes: []byte(")"),
			}
			newTokens = append(newTokens, fnEnd)
			body.SetAttributeRaw(attributeName, newTokens)
		}
	}
}
