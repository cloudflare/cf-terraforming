resource "cloudflare_zero_trust_risk_scoring_integration" "example_zero_trust_risk_scoring_integration" {
  account_id = "account_id"
  integration_type = "Okta"
  tenant_url = "https://example.com"
  reference_id = "reference_id"
}
