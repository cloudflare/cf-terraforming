resource "cloudflare_zero_trust_tunnel_cloudflared" "terraform_managed_resource" {
  account_id = "699d98642c564d2e855e9661899b7252"
  name = "blog"
  config_src = "local"
  tunnel_secret = "AQIDBAUGBwgBAgMEBQYHCAECAwQFBgcIAQIDBAUGBwg="
}
