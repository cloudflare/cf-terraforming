resource "cloudflare_email_security_trusted_domains" "terraform_managed_resource" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  is_recent = true
  is_regex = false
  is_similarity = false
  pattern = "example.com"
  comments = null
}
