resource "cloudflare_ruleset" "terraform_managed_resource" {
  kind    = "zone"
  name    = "default"
  phase   = "http_request_cache_settings"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  rules = [{
    action = "set_cache_settings"
    action_parameters = {
      browser_ttl = {
        mode = "respect_origin"
      }
      cache = true
      cache_key = {
        cache_by_device_type  = true
        cache_deception_armor = true
        custom_key = {
          host = {
            resolved = false
          }
          query_string = {
            exclude = ["*"]
          }
        }
        ignore_query_strings_order = false
      }
      edge_ttl = {
        default = 30
        mode    = "override_origin"
        status_code_ttl = [{
          status_code = 100
          value       = 30
          }, {
          status_code_range = {
            from = 100
            to   = 106
          }
          value = 5
          }, {
          status_code_range = {
            from = 130
            to   = 162
          }
          value = 31536000
        }]
      }
      origin_error_page_passthru = true
      respect_strong_etags       = true
      serve_stale = {
        disable_stale_while_updating = true
      }
    }
    description  = "test cache rule"
    enabled      = false
    expression   = "(http.host eq \"example.com\")"
    id           = null
    last_updated = "2022-09-28T17:21:21.510301Z"
    ref          = "0f24aab3002347a9a4ac01520e6893d0"
    version      = "3"
    }, {
    action = "set_cache_settings"
    action_parameters = {
      cache = false
      edge_ttl = {
        default = 60
        mode    = "override_origin"
      }
    }
    description  = "/status/202"
    enabled      = true
    expression   = "(http.host eq \"example.com\")"
    id           = null
    last_updated = "2022-09-21T16:36:00.999083Z"
    ref          = "e5f1bd1386b4464aa8d726ba1e0d51ad"
    version      = "2"
  }]
}

