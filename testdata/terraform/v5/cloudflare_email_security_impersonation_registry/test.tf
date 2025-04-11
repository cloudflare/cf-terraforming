resource "cloudflare_email_security_impersonation_registry" "terraform_managed_resource" {
  account_id     = "f037e56e89293a057740de681ac9abbe"
  email          = "email"
  is_email_regex = true
  name           = "name"
}

