resource "cloudflare_ruleset" "terraform_managed_resource" {
  kind    = "zone"
  name    = "zone"
  phase   = "http_request_firewall_managed"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  rules {
    action = "execute"
    action_parameters {
      id = "efb7b8c949ac4650a09736fc376e9aee"
      overrides {
        action = "log"
        categories {
          category = "paranoia-level-2"
          enabled  = false
        }
        enabled = false
        rules {
          enabled = false
          id      = "6179ae15870a4bb7b2d480d4843b323c"
        }
      }
    }
    description = "zone"
    enabled     = false
    expression  = "(http.cookie eq \"jb_testing=true\")"
  }
}
