resource "cloudflare_ruleset" "terraform_managed_resource" {
  kind    = "zone"
  name    = "default"
  phase   = "http_request_cache_settings"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  rules {
    action      = "set_cache_settings"
    description = "test cache rule"
    enabled     = false
    expression  = "(http.host eq \"example.com\")"
    action_parameters {
      browser_ttl {
        mode = "respect_origin"
      }
      cache_key {
        custom_key {
          host {
            resolved = false
          }
          query_string {
            exclude = ["*"]
          }
        }
        cache_by_device_type       = true
        cache_deception_armor      = true
        ignore_query_strings_order = false
      }
      edge_ttl {
        status_code_ttl {
          status_code = 100
          value       = 30
        }
        default = 30
        mode    = "override_origin"
      }
      serve_stale {
        disable_stale_while_updating = true
      }
      cache                      = true
      origin_error_page_passthru = true
      respect_strong_etags       = true
    }
  }
  rules {
    action      = "set_cache_settings"
    description = "/status/202"
    enabled     = true
    expression  = "(http.host eq \"example.com\")"
    action_parameters {
      edge_ttl {
        default = 60
        mode    = "override_origin"
      }
      cache = false
    }
  }
}
