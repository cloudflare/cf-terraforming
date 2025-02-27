resource "cloudflare_bot_management" "terraform_managed_resource" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  ai_bots_protection = "block"
  enable_js = true
  fight_mode = true
}
