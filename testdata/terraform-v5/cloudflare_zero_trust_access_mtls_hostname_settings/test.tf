resource "cloudflare_zero_trust_access_mtls_hostname_settings" "terraform_managed_resource" {
  settings = [{
    china_network = false
    client_certificate_forwarding = true
    hostname = "admin.example.com"
  }]
  zone_id = "zone_id"
}
