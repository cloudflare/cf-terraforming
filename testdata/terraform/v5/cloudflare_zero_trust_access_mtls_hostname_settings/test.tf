resource "cloudflare_zero_trust_access_mtls_hostname_settings" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  settings = [{
    china_network                 = false
    client_certificate_forwarding = true
    hostname                      = "terraform.cfapi.net"
  }]
}

