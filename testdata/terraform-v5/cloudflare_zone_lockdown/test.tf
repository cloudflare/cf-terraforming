resource "cloudflare_zone_lockdown" "terraform_managed_resource" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  configurations = [{
    target = "ip"
    value = "198.51.100.4"
  }]
  urls = ["shop.example.com/*"]
}
