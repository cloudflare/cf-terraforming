resource "cloudflare_zero_trust_risk_scoring_integration" "terraform_managed_resource" {
  account_id       = "f037e56e89293a057740de681ac9abbe"
  active           = true
  integration_type = "Okta"
  reference_id     = "3701e70b-5de4-4d9d-8972-f5aa994d9057"
  tenant_url       = "https://example.oktapreview.com"
}

