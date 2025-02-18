resource "cloudflare_spectrum_application" "example_spectrum_application" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  dns = {
    name = "ssh.example.com"
    type = "CNAME"
  }
  ip_firewall = true
  protocol = "tcp/22"
  proxy_protocol = "off"
  tls = "off"
  traffic_type = "direct"
  argo_smart_routing = true
  edge_ips = {
    connectivity = "all"
    type = "dynamic"
  }
  origin_direct = ["tcp://127.0.0.1:8080"]
  origin_dns = {
    name = "origin.example.com"
    ttl = 600
    type = ""
  }
  origin_port = 22
}
