resource "cloudflare_ruleset" "terraform_managed_resource" {
  kind    = "zone"
  name    = "default"
  phase   = "http_config_settings"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  rules {
    action = "set_config"
    action_parameters {
      content_converter = false
    }
    description = "convert HTML to Markdown"
    enabled     = true
    expression  = "true"
  }
}
