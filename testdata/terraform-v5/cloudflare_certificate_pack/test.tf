resource "cloudflare_certificate_pack" "example_certificate_pack" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  certificate_authority = "google"
  hosts = ["example.com", "*.example.com", "www.example.com"]
  type = "advanced"
  validation_method = "txt"
  validity_days = 14
  cloudflare_branding = false
}
