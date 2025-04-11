resource "cloudflare_email_security_trusted_domains" "terraform_managed_resource" {
  account_id    = "f037e56e89293a057740de681ac9abbe"
  is_recent     = true
  is_regex      = false
  is_similarity = false
  pattern       = "example.com"
}

