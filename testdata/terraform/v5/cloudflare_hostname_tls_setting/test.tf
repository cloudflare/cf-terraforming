resource "cloudflare_hostname_tls_setting" "terraform_managed_resource_0" {
  hostname   = "cdvrjwgmzv.example.com"
  setting_id = "ciphers"
  zone_id    = "0da42c8d2132a9ddaf714f9e7c920711"
  value      = ["AES128-SHA", "ECDHE-RSA-AES256-SHA"]
}

resource "cloudflare_hostname_tls_setting" "terraform_managed_resource_1" {
  hostname   = "example.com"
  setting_id = "ciphers"
  zone_id    = "0da42c8d2132a9ddaf714f9e7c920711"
  value      = ["ECDHE-RSA-AES128-GCM-SHA256", "AES128-GCM-SHA256"]
}

