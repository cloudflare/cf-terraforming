resource "cloudflare_zero_trust_gateway_proxy_endpoint" "example_zero_trust_gateway_proxy_endpoint" {
  account_id = "699d98642c564d2e855e9661899b7252"
  ips = ["192.0.2.1/32"]
  name = "Devops team"
}
