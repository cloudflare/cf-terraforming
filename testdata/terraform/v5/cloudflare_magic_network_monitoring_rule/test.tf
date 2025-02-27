resource "cloudflare_magic_network_monitoring_rule" "terraform_managed_resource" {
  account_id = "6f91088a406011ed95aed352566e8d4c"
  duration = "300s"
  name = "my_rule_1"
  automatic_advertisement = true
  bandwidth = 1000
  packet_threshold = 10000
  prefixes = ["203.0.113.1/32"]
}
