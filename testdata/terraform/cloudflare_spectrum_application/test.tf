resource "cloudflare_spectrum_application" "terraform_managed_resource" {
  argo_smart_routing = true
  edge_ips           = ["198.51.100.1"]
  ip_firewall        = true
  origin_direct      = ["tcp://192.0.2.1:22"]
  protocol           = "tcp/22"
  proxy_protocol     = "off"
  tls                = "full"
  traffic_type       = "direct"
  zone_id            = "0da42c8d2132a9ddaf714f9e7c920711"
  dns {
    name = "ssh.example.com"
    type = "CNAME"
  }
}
