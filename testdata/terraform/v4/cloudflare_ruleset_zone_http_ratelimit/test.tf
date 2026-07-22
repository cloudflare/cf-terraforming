resource "cloudflare_ruleset" "terraform_managed_resource" {
  kind    = "zone"
  name    = "default"
  phase   = "http_ratelimit"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  rules {
    action      = "block"
    description = "fwewe"
    enabled     = false
    expression  = "(http.cookie eq \"namwe=value\")"
    ratelimit {
      characteristics     = ["ip.src", "cf.colo.id"]
      mitigation_timeout  = 30
      period              = 60
      requests_per_period = 100
      requests_to_origin  = true
    }
  }
}
