resource "cloudflare_ruleset" "terraform_managed_resource" {
  kind    = "zone"
  name    = "default"
  phase   = "http_request_dynamic_redirect"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  rules = [{
    action = "redirect"
    action_parameters = {
      from_value = {
        preserve_query_string = true
        status_code           = 301
        target_url = {
          value = "https://example.com/foo"
        }
      }
    }
    enabled      = true
    expression   = "true"
    id           = null
    last_updated = "2025-03-25T22:38:12.845519Z"
    ref          = "jacob1"
    version      = "2"
    }, {
    action = "redirect"
    action_parameters = {
      from_value = {
        preserve_query_string = true
        status_code           = 301
        target_url = {
          value = "https://example.com/example1"
        }
      }
    }
    enabled      = true
    expression   = "true"
    id           = null
    last_updated = "2025-03-25T22:38:12.845519Z"
    ref          = "jacob2"
    version      = "3"
  }]
}

