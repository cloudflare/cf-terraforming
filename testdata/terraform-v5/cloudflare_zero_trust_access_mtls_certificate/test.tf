resource "cloudflare_zero_trust_access_mtls_certificate" "example_zero_trust_access_mtls_certificate" {
  certificate = <<EOT
  -----BEGIN CERTIFICATE-----
  MIIGAjCCA+qgAwIBAgIJAI7kymlF7CWT...N4RI7KKB7nikiuUf8vhULKy5IX10
  DrUtmu/B
  -----END CERTIFICATE-----
  EOT
  name = "Allow devs"
  zone_id = "zone_id"
  associated_hostnames = ["admin.example.com"]
}
