resource "cloudflare_zero_trust_dlp_custom_profile" "terraform_managed_resource" {
  account_id           = "f037e56e89293a057740de681ac9abbe"
  ai_context_enabled   = false
  allowed_match_count  = 0
  confidence_threshold = "low"
  description          = "custom profile"
  name                 = "psuhmwlpqf"
  ocr_enabled          = true
  context_awareness = {
    enabled = false
    skip = {
      files = false
    }
  }
  entries = [{
    created_at = "2024-05-15T06:02:05Z"
    enabled    = true
    id         = "34f2bd4b-5069-4f5b-a22e-3f7878912032"
    name       = "psuhmwlpqf_entry1"
    pattern = {
      regex      = "^4[0-9]"
      validation = "luhn"
    }
    profile_id = "38f45ad8-476e-4b56-ad16-42f364250802"
    type       = "custom"
    updated_at = "2024-05-15T06:02:05Z"
  }]
}

