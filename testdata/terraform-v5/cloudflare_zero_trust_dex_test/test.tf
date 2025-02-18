resource "cloudflare_zero_trust_dex_test" "example_zero_trust_dex_test" {
  account_id = "699d98642c564d2e855e9661899b7252"
  data = {
    host = "https://dash.cloudflare.com"
    kind = "http"
    method = "GET"
  }
  enabled = true
  interval = "30m"
  name = "HTTP dash health check"
  description = "Checks the dash endpoint every 30 minutes"
  target_policies = [{
    id = "id"
    default = true
    name = "name"
  }]
  targeted = true
}
