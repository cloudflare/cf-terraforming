resource "cloudflare_email_security_block_sender" "terraform_managed_resource" {
  account_id   = "f037e56e89293a057740de681ac9abbe"
  comments     = "comments"
  is_regex     = true
  pattern      = "x"
  pattern_type = "EMAIL"
}

