resource "cloudflare_zone_lockdown" "terraform_managed_resource" {
  description = "Restrict access to these endpoints to requests from a known IP address"
  paused = false
  urls = [ "api.mysite.com/some/endpoint*" ]
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  configurations {
    target = "ip"
    value = "198.51.100.4"
  }
}
