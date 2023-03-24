resource "cloudflare_ruleset" "terraform_managed_resource" {
  kind    = "zone"
  name    = "default"
  phase   = "http_request_firewall_custom"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  rules {
    action = "skip"
    action_parameters {
      phases   = ["http_ratelimit", "http_request_firewall_managed"]
      products = ["zoneLockdown", "uaBlock", "bic", "hot", "securityLevel", "rateLimit", "waf"]
      ruleset  = "current"
    }
    description  = "test.example.com"
    enabled      = true
    expression   = "(http.host eq \"test.example.com\")"
    id           = "88dcb30401e348ba9e1352c2598f2a4c"
    last_updated = "2022-11-24T14:24:14.756247Z"
    logging {
      enabled = true
    }
    ref     = "88dcb30401e348ba9e1352c2598f2a4c"
    version = "2"
  }
  rules {
    action       = "challenge"
    description  = "customRule-test"
    enabled      = true
    expression   = "(cf.bot_management.score eq 50 and cf.bot_management.static_resource)"
    id           = "b3cc5e4cc6604f9d90a6a106df867760"
    last_updated = "2022-11-07T19:03:05.198191Z"
    ref          = "b3cc5e4cc6604f9d90a6a106df867760"
    version      = "29"
  }
  rules {
    action       = "log"
    description  = "AWAF ML"
    enabled      = false
    expression   = "(cf.waf.score le 20)"
    id           = "1ecf73bdf7bd4227969a734412b13ad1"
    last_updated = "2022-12-09T16:53:19.003821Z"
    ref          = "1ecf73bdf7bd4227969a734412b13ad1"
    version      = "7"
  }
}
