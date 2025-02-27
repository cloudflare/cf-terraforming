resource "cloudflare_record" "terraform_managed_resource" {
  content = "255.2.0.192.in-addr.arpa"
  name    = "example.com"
  proxied = false
  ttl     = 1
  type    = "PTR"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
}
