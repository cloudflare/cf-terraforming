resource "cloudflare_web_analytics_site" "terraform_managed_resource" {
  account_id   = "f037e56e89293a057740de681ac9abbe"
  auto_install = true
  enabled      = true
  lite         = false
  zone_tag     = "0da42c8d2132a9ddaf714f9e7c920711"
}

