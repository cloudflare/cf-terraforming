resource "cloudflare_zone_lockdown" "terraform_managed_resource" {
  paused  = false
  urls    = ["ephnewtbrw.terraform.cfapi.net/*"]
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  configurations = [{
    target = "ip"
    value  = "198.51.100.4"
  }]
}

