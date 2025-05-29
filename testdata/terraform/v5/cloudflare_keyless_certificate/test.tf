resource "cloudflare_keyless_certificate" "terraform_managed_resource" {
  certificate = "-----INSERT CERTIFICATE-----"
  enabled     = false
  host        = "terraform.cfapi.net"
  name        = "ydvgqcbbcq"
  port        = 24008
  zone_id     = "0da42c8d2132a9ddaf714f9e7c920711"
}

