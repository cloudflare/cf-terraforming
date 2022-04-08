
resource "cloudflare_record" "terraform_managed_resource" {
  name    = "txtspf.example.com"
  proxied = false
  ttl     = 1
  type    = "TXT"
  value   = "\"v=spf1 include:%%{ir}.%%{v}.%%{d}.spf.has.pphosted.com include:amazonses.com ~all\""
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
}