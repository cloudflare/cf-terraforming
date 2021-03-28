resource "cloudflare_record" "terraform_managed_resource" {
  value = "198.51.100.4"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  name = "example.com"
  data = {
  }
  proxied = false
  type = "A"
  created_on = "2014-01-01T05:20:00.12345Z"
  proxiable = true
  ttl = 120
  modified_on = "2014-01-01T05:20:00.12345Z"
}
