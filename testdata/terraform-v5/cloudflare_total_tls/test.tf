resource "cloudflare_total_tls" "example_total_tls" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  enabled = true
  certificate_authority = "google"
}
