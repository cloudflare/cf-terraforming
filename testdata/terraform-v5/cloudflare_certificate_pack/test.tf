resource "cloudflare_certificate_pack" "terraform_managed_resource" {
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  certificate_authority = "digicert"
  hosts = ["example.com", "*.example.com", "www.example.com"]
  type = "advanced"
  validation_method = "txt"
  validity_days = 365
  cloudflare_branding = false
}
