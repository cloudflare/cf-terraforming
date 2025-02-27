resource "cloudflare_total_tls" "terraform_managed_resource" {
  certificate_authority = "google"
  enabled               = true
  zone_id               = "0da42c8d2132a9ddaf714f9e7c920711"
}
