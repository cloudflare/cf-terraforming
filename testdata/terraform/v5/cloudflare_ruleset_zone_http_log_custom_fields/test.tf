resource "cloudflare_ruleset" "terraform_managed_resource" {
  kind    = "zone"
  name    = "zone"
  phase   = "http_log_custom_fields"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  rules = [{
    action = "log_custom_field"
    action_parameters = {
      cookie_fields = [{
        name = "cookie"
        }, {
        name = "fields"
      }]
      request_fields = [{
        name = "request"
        }, {
        name = "fields"
      }]
      response_fields = [{
        name = "response"
        }, {
        name = "fields"
      }]
    }
    description  = "zone"
    enabled      = true
    expression   = "true"
    id           = null
    last_updated = "2022-07-22T12:34:45.479429Z"
    ref          = "17a0d1e23a3444ccbd5e58fc7793649a"
    version      = "1"
  }]
}

