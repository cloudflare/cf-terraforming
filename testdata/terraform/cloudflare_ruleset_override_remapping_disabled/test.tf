resource "cloudflare_ruleset" "terraform_managed_resource" {
  kind    = "zone"
  name    = "zone"
  phase   = "http_request_firewall_managed"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  rules {
    action      = "execute"
    description = "zone"
    enabled     = false
    expression  = "(http.cookie eq \"jb_testing=true\")"
    action_parameters {
      overrides {
        action = "log"
        status = "disabled"
      }
      id      = "efb7b8c949ac4650a09736fc376e9aee"
      version = "latest"
    }
  }
}
