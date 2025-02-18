resource "cloudflare_turnstile_widget" "example_turnstile_widget" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  domains = ["203.0.113.1", "cloudflare.com", "blog.example.com"]
  mode = "non-interactive"
  name = "blog.cloudflare.com login form"
  bot_fight_mode = false
  clearance_level = "no_clearance"
  ephemeral_id = false
  offlabel = false
  region = "world"
}
