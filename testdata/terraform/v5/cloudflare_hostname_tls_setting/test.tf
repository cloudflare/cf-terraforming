resource "cloudflare_hostname_tls_setting" "terraform_managed_resource" {
  zone_id    = "0da42c8d2132a9ddaf714f9e7c920711"
  hostname   = "terraform.cfapi.net"
  setting_id = "ciphers"
  value      = ["ECDHE-RSA-AES128-GCM-SHA256", "AES128-GCM-SHA256"]
}
