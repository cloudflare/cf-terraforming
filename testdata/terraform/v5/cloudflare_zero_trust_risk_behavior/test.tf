resource "cloudflare_zero_trust_risk_behavior" "terraform_managed_resource" {
  account_id = "account_id"
  behaviors = {
    foo = {
      enabled = true
      risk_level = "low"
    }
  }
}
