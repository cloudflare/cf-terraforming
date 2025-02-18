resource "cloudflare_zone_lockdown" "example_zone_lockdown" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  configurations = [{
    target = "ip"
    value = "198.51.100.4"
  }]
  urls = ["shop.example.com/*"]
}
