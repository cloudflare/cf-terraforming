resource "cloudflare_zero_trust_tunnel_cloudflared_route" "example_zero_trust_tunnel_cloudflared_route" {
  account_id = "699d98642c564d2e855e9661899b7252"
  network = "172.16.0.0/16"
  tunnel_id = "f70ff985-a4ef-4643-bbbc-4a0ed4fc8415"
  comment = "Example comment for this route."
  virtual_network_id = "f70ff985-a4ef-4643-bbbc-4a0ed4fc8415"
}
