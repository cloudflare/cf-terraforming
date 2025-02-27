resource "cloudflare_user_agent_blocking_rule" "terraform_managed_resource" {
  description = "My description 1"
  mode        = "js_challenge"
  paused      = false
  zone_id     = "0da42c8d2132a9ddaf714f9e7c920711"
  configuration {
    target = "ua"
    value  = "Chrome"
  }
}
