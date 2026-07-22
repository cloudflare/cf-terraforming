resource "cloudflare_record" "terraform_managed_resource" {
  content = "\"v=spf1 include:%%{ir}.%%{v}.%%{d}.spf.has.pphosted.com include:amazonses.com ~all\""
  name    = "txtspf"
  proxied = false
  ttl     = 1
  type    = "TXT"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
}
