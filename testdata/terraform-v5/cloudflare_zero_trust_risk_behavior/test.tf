resource "cloudflare_zero_trust_risk_behavior" "example_zero_trust_risk_behavior" {
  account_id = "account_id"
  behaviors = {
    foo = {
      enabled = true
      risk_level = "low"
    }
  }
}
