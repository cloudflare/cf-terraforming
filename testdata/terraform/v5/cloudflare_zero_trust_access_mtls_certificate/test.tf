resource "cloudflare_zero_trust_access_mtls_certificate" "terraform_managed_resource" {
  account_id           = "f037e56e89293a057740de681ac9abbe"
  associated_hostnames = []
  certificate          = "-----INSERT CERTIFICATE-----"
  name                 = "Allow devs"
}

