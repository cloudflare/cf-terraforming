resource "cloudflare_zero_trust_dlp_entry" "example_zero_trust_dlp_entry" {
  account_id = "account_id"
  enabled = true
  name = "name"
  pattern = {
    regex = "regex"
    validation = "luhn"
  }
  profile_id = "182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"
}
