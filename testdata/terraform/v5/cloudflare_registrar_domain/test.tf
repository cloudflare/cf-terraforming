resource "cloudflare_registrar_domain" "terraform_managed_resource_0" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  auto_renew  = true
  domain_name = "12345678dnstest.org"
  locked      = true
  privacy     = true
}

resource "cloudflare_registrar_domain" "terraform_managed_resource_1" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  auto_renew  = true
  domain_name = "1234test.dev"
  locked      = true
  privacy     = true
}

