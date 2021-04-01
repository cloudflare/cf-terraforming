resource "cloudflare_record" "terraform_managed_resource" {
  created_on = "2014-01-01T05:20:00.12345Z"
  modified_on = "2014-01-01T05:20:00.12345Z"
  name = "example.com"
  proxiable = true
  proxied = false
  ttl = 120
  type = "A"
  value = "198.51.100.4"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
}
