resource "cloudflare_bot_management" "terraform_managed_resource" {
  ai_bots_protection              = "block"
  enable_js                       = true
  optimize_wordpress              = true
  sbfm_definitely_automated       = "block"
  sbfm_likely_automated           = "managed_challenge"
  sbfm_static_resource_protection = false
  sbfm_verified_bots              = "allow"
  zone_id                         = "0da42c8d2132a9ddaf714f9e7c920711"
}
