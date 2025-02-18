resource "cloudflare_hostname_tls_setting" "example_hostname_tls_setting" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  setting_id = "ciphers"
  hostname = "app.example.com"
  value = ["ECDHE-RSA-AES128-GCM-SHA256", "AES128-GCM-SHA256"]
}
