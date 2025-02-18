resource "cloudflare_magic_wan_ipsec_tunnel" "example_magic_wan_ipsec_tunnel" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  cloudflare_endpoint = "203.0.113.1"
  interface_address = "192.0.2.0/31"
  name = "IPsec_1"
  customer_endpoint = "203.0.113.1"
  description = "Tunnel for ISP X"
  health_check = {
    direction = "unidirectional"
    enabled = true
    rate = "low"
    target = {
      saved = "203.0.113.1"
    }
    type = "reply"
  }
  psk = "O3bwKSjnaoCxDoUxjcq4Rk8ZKkezQUiy"
  replay_protection = false
}
