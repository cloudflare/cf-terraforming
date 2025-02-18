resource "cloudflare_zero_trust_dlp_custom_profile" "example_zero_trust_dlp_custom_profile" {
  account_id = "account_id"
  profiles = [{
    entries = [{
      enabled = true
      name = "name"
      pattern = {
        regex = "regex"
        validation = "luhn"
      }
    }]
    name = "name"
    allowed_match_count = 5
    confidence_threshold = "confidence_threshold"
    context_awareness = {
      enabled = true
      skip = {
        files = true
      }
    }
    description = "description"
    ocr_enabled = true
    shared_entries = [{
      enabled = true
      entry_id = "182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"
      entry_type = "custom"
    }]
  }]
}
