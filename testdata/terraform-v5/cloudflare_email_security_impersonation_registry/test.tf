resource "cloudflare_email_security_impersonation_registry" "example_email_security_impersonation_registry" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  email = "email"
  is_email_regex = true
  name = "name"
}
