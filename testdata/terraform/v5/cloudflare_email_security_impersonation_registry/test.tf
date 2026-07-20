resource "cloudflare_email_security_impersonation_registry" "terraform_managed_resource" {
  account_id                 = "f037e56e89293a057740de681ac9abbe"
  directory_id               = 0
  email                      = "email"
  external_directory_node_id = "12444788"
  is_email_regex             = true
  name                       = "name"
  provenance                 = "A1S_INTERNAL"
}

