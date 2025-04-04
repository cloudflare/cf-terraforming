resource "cloudflare_zero_trust_device_posture_rule" "terraform_managed_resource_0" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  description = "My description"
  name        = "wfbhpishtq"
  schedule    = "24h"
  type        = "disk_encryption"
  input = {
    checkDisks = []
    requireAll = true
  }
  match = [{
    platform = "mac"
  }]
}

resource "cloudflare_zero_trust_device_posture_rule" "terraform_managed_resource_1" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  description = "My description"
  name        = "jsdpiddjgv"
  schedule    = "24h"
  type        = "disk_encryption"
  input = {
    checkDisks = []
    requireAll = true
  }
  match = [{
    platform = "mac"
  }]
}

