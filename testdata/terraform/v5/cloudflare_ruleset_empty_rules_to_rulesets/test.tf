resource "cloudflare_ruleset" "terraform_managed_resource" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  description = "Account-level WAF configuration"
  kind        = "root"
  name        = "Account WAF Ruleset"
  phase       = "http_request_firewall_managed"
  rules = [{
    action = "skip"
    action_parameters = {
      rulesets = ["4814384a9e5d4991b9815dcfc25d2f1f", "c2e184081120413c86c3ab7e14069605"]
    }
    description  = "Skip all managed rules for API subdomain"
    enabled      = true
    expression   = "(http.host eq \"api.example.com\")"
    id           = null
    last_updated = "2024-01-15T10:30:00.123456Z"
    logging = {
      enabled = true
    }
    ref     = "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
    version = "1"
  }]
}
