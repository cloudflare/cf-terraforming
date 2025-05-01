resource "cloudflare_web_analytics_rule" "terraform_managed_resource_0" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  host       = "*"
  inclusive  = true
  is_paused  = false
  paths      = ["*"]
  ruleset_id = "2fa89d8f-35f7-49ef-87d3-f24e866a5d5e"
}

resource "cloudflare_web_analytics_rule" "terraform_managed_resource_1" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  host       = "*.example.com"
  inclusive  = false
  is_paused  = false
  paths      = ["v1/images/*"]
  ruleset_id = "2fa89d8f-35f7-49ef-87d3-f24e866a5d5e"
}

