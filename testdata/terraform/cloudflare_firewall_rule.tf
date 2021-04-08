resource "cloudflare_firewall_rule" "terraform_managed_resource" {
  action = "block"
  description = "Blocks traffic identified during investigation for MIR-31"
  filter_id = "372e67954025e0ba6aaa6d586b9e0b61"
  paused = false
  priority = 50
  products = [ "waf" ]
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
}
