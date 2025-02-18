resource "cloudflare_zero_trust_tunnel_cloudflared_virtual_network" "example_zero_trust_tunnel_cloudflared_virtual_network" {
  account_id = "699d98642c564d2e855e9661899b7252"
  name = "us-east-1-vpc"
  comment = "Staging VPC for data science"
  is_default = true
}
