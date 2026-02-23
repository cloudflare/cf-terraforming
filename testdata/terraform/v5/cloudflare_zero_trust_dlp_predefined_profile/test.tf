resource "cloudflare_zero_trust_dlp_predefined_profile" "terraform_managed_resource_0" {
  account_id          = "f037e56e89293a057740de681ac9abbe"
  ai_context_enabled  = false
  allowed_match_count = 0
  ocr_enabled         = false
  profile_id          = "c8932cc4-3312-4152-8041-f3f257122dc4"
  entries = [{
    confidence = {
      ai_context_available = true
      available            = false
    }
    enabled    = false
    id         = "d8fcfc9c-773c-405e-8426-21ecbb67ba93"
    name       = "Amazon AWS Access Key ID"
    profile_id = "c8932cc4-3312-4152-8041-f3f257122dc4"
    type       = "predefined"
    }, {
    confidence = {
      ai_context_available = false
      available            = false
    }
    enabled    = false
    id         = "2c0e33e1-71da-40c8-aad3-32e674ad3d96"
    name       = "Amazon AWS Secret Access Key"
    profile_id = "c8932cc4-3312-4152-8041-f3f257122dc4"
    type       = "predefined"
    }, {
    confidence = {
      ai_context_available = true
      available            = false
    }
    enabled    = false
    id         = "6c6579e4-d832-42d5-905c-8e53340930f2"
    name       = "Google GCP API Key"
    profile_id = "c8932cc4-3312-4152-8041-f3f257122dc4"
    type       = "predefined"
    }, {
    confidence = {
      ai_context_available = true
      available            = false
    }
    enabled    = false
    id         = "4e92c006-3802-4dff-bbe1-8e1513b1c92a"
    name       = "Microsoft Azure Client Secret"
    profile_id = "c8932cc4-3312-4152-8041-f3f257122dc4"
    type       = "predefined"
    }, {
    confidence = {
      ai_context_available = false
      available            = false
    }
    enabled    = false
    id         = "5c713294-2375-4904-abcf-e4a15be4d592"
    name       = "SSH Private Key"
    profile_id = "c8932cc4-3312-4152-8041-f3f257122dc4"
    type       = "predefined"
  }]
}

resource "cloudflare_zero_trust_dlp_predefined_profile" "terraform_managed_resource_1" {
  account_id          = "f037e56e89293a057740de681ac9abbe"
  ai_context_enabled  = false
  allowed_match_count = 0
  ocr_enabled         = false
  profile_id          = "56a8c060-01bb-4f89-ba1e-3ad42770a342"
  entries = [{
    confidence = {
      ai_context_available = true
      available            = false
    }
    enabled    = false
    id         = "d8fcfc9c-773c-405e-8426-21ecbb67ba93"
    name       = "Amazon AWS Access Key ID"
    profile_id = "c8932cc4-3312-4152-8041-f3f257122dc4"
    type       = "predefined"
    }, {
    confidence = {
      ai_context_available = false
      available            = false
    }
    enabled    = false
    id         = "2c0e33e1-71da-40c8-aad3-32e674ad3d96"
    name       = "Amazon AWS Secret Access Key"
    profile_id = "c8932cc4-3312-4152-8041-f3f257122dc4"
    type       = "predefined"
    }, {
    confidence = {
      ai_context_available = true
      available            = false
    }
    enabled    = false
    id         = "6c6579e4-d832-42d5-905c-8e53340930f2"
    name       = "Google GCP API Key"
    profile_id = "c8932cc4-3312-4152-8041-f3f257122dc4"
    type       = "predefined"
    }, {
    confidence = {
      ai_context_available = true
      available            = false
    }
    enabled    = false
    id         = "4e92c006-3802-4dff-bbe1-8e1513b1c92a"
    name       = "Microsoft Azure Client Secret"
    profile_id = "c8932cc4-3312-4152-8041-f3f257122dc4"
    type       = "predefined"
    }, {
    confidence = {
      ai_context_available = false
      available            = false
    }
    enabled    = false
    id         = "5c713294-2375-4904-abcf-e4a15be4d592"
    name       = "SSH Private Key"
    profile_id = "c8932cc4-3312-4152-8041-f3f257122dc4"
    type       = "predefined"
  }]
}

