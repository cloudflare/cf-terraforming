resource "cloudflare_record" "terraform_managed_resource" {
  content = "198.51.100.4"
  name    = "subdomain"
  proxied = false
  ttl     = 120
  type    = "A"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
}
