resource "cloudflare_magic_network_monitoring_configuration" "example_magic_network_monitoring_configuration" {
  account_id = "6f91088a406011ed95aed352566e8d4c"
  default_sampling = 1
  name = "cloudflare user\'s account"
  router_ips = ["203.0.113.1"]
  warp_devices = [{
    id = "5360368d-b351-4791-abe1-93550dabd351"
    name = "My warp device"
    router_ip = "203.0.113.1"
  }]
}
