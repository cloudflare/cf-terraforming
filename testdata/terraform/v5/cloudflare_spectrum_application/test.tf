resource "cloudflare_spectrum_application" "terraform_managed_resource" {
  ip_firewall    = false
  origin_direct  = ["tcp://128.66.0.4:3389"]
  protocol       = "rdp"
  proxy_protocol = "off"
  tls            = "off"
  traffic_type   = "direct"
  zone_id        = "0da42c8d2132a9ddaf714f9e7c920711"
  dns = {
    name = "ledfaootii.terraform.cfapi.net"
    type = "CNAME"
  }
  edge_ips = {
    connectivity = "all"
    type         = "dynamic"
  }
}

