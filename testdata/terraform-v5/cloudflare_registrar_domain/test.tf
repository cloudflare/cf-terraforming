resource "cloudflare_registrar_domain" "terraform_managed_resource" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  domain_name = "cloudflare.com"
  auto_renew = true
  locked = false
  privacy = true
}
