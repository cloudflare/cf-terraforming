resource "cloudflare_record" "terraform_managed_resource" {
  data = {
    flags = 0
    tag   = "issuewild"
    value = "example.com"
  }
  name    = "example.com"
  proxied = false
  ttl     = 120
  type    = "CAA"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
}
