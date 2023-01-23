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
      cache                      = true
      origin_error_page_passthru = true
      respect_strong_etags       = true
      browser_ttl {
        mode = "respect_origin"
      }
      cache_key {
        cache_by_device_type       = true
        cache_deception_armor      = true
        ignore_query_strings_order = false
        custom_key {
          host {
            resolved = false
          }
          query_string {
            exclude = ["*"]
          }
        }
      }
      edge_ttl {
        default = 30
        mode    = "override_origin"
        status_code_ttl {
          value = 1
          status_code_range {
            from = 101
            to   = 103
          }
        }
        status_code_ttl {
          value = 1
          status_code_range {
            from = 106
            to   = 110
          }
        }
      }
      serve_stale {
        disable_stale_while_updating = true
      }
    }
  }
  rules {
    action      = "set_cache_settings"
    description = "test cache rule 2"
    enabled     = true
    expression  = "(http.host eq \"example.com\")"
    action_parameters {
      cache = true
      edge_ttl {
        mode = "respect_origin"
        status_code_ttl {
          value = 10
          status_code_range {
            from = 1
            to   = 2
          }
        }
        status_code_ttl {
          value = 1
          status_code_range {
            from = 3
            to   = 4
          }
        }
      }
    }
  }
  rules {
    action      = "set_cache_settings"
    description = "test cache rule"
    enabled     = false
    expression  = "(http.host eq \"example.com\")"
    action_parameters {
      cache                = true
      respect_strong_etags = true
      browser_ttl {
        mode = "respect_origin"
      }
      cache_key {
        cache_by_device_type = true
        custom_key {
          host {
            resolved = false
          }
          query_string {
            include = ["*"]
          }
        }
      }
      edge_ttl {
        default = 1
        mode    = "override_origin"
        status_code_ttl {
          status_code = 100
          value       = 5
        }
      }
      serve_stale {
        disable_stale_while_updating = true
      }
    }
  }
}
