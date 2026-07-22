resource "cloudflare_snippet_rules" "terraform_managed_resource" {
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  rules = [{
    description  = ""
    enabled      = true
    expression   = "(http.request.full_uri wildcard \"/hello\")"
    snippet_name = "remove_query_strings_template"
  }]
}

