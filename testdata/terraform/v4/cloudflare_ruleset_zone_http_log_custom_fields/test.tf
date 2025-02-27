resource "cloudflare_ruleset" "terraform_managed_resource" {
  kind    = "zone"
  name    = "zone"
  phase   = "http_log_custom_fields"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  rules {
    action = "log_custom_field"
    action_parameters {
      cookie_fields   = ["cookie", "fields"]
      request_fields  = ["request", "fields"]
      response_fields = ["response", "fields"]
    }
    description = "zone"
    enabled     = true
    expression  = "true"
  }
}
