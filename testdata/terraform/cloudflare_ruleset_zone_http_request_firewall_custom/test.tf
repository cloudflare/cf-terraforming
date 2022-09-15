resource "cloudflare_ruleset" "terraform_managed_resource" {
  kind    = "zone"
  name    = "default"
  phase   = "http_request_firewall_custom"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  rules {
    action      = "skip"
    description = "firewall rule"
    enabled     = true
    expression  = "(http.request.uri.path contains \"/filters\")"
    action_parameters {
      rules = {
        "efb7b8c949ac4650a09736fc376e9aee" = "062a7840e0cb47f7b36acd2d507ce584,5cLhGXtTafjwPkdy8fmW5QvPiokBuZhi"
      }
    }
    logging {
      status = "enabled"
    }
  }
  rules {
    action      = "skip"
    description = "test skip rule on ip "
    enabled     = true
    expression  = "(ip.src eq 1.2.3.4)"
    action_parameters {
      ruleset = "current"
    }
    logging {
      status = "disabled"
    }
  }
}
