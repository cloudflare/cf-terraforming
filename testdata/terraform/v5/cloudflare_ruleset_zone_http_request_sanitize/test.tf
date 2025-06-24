resource "cloudflare_ruleset" "terraform_managed_resource" {
  kind    = "zone"
  name    = "Zone sanitize ruleset"
  phase   = "http_request_sanitize"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  rules = [{
    action = "execute"
    action_parameters = {
      id = "70339d97bdb34195bbf054b1ebe81f76"
      overrides = {
        rules = [{
          enabled = true
          id      = "78723a9e0c7c4c6dbec5684cb766231d"
        }]
      }
      version = "latest"
    }
    enabled      = true
    expression   = "true"
    id           = null
    last_updated = "2021-08-19T23:41:34.985519Z"
    ref          = "0789dc4343054d1e981f8c44bedc6fbd"
    version      = "1"
  }]
}

