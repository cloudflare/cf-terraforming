resource "cloudflare_email_security_block_sender" "example_email_security_block_sender" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  is_regex = false
  pattern = "test@example.com"
  pattern_type = "EMAIL"
  comments = "block sender with email test@example.com"
}
