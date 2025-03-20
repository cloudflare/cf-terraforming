resource "cloudflare_bot_management" "terraform_managed_resource" {
  ai_bots_protection              = "disabled"
  enable_js                       = true
  optimize_wordpress              = true
  sbfm_definitely_automated       = "managed_challenge"
  sbfm_likely_automated           = "block"
  sbfm_static_resource_protection = false
  sbfm_verified_bots              = "allow"
  suppress_session_score          = false
  zone_id                         = "0da42c8d2132a9ddaf714f9e7c920711"
}

