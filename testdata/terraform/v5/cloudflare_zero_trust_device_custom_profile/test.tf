resource "cloudflare_zero_trust_device_custom_profile" "terraform_managed_resource_0" {
  account_id            = "f037e56e89293a057740de681ac9abbe"
  allow_mode_switch     = true
  allow_updates         = true
  allowed_to_leave      = true
  auto_connect          = 0
  captive_portal        = 5
  description           = "xocmddmeyz"
  disable_auto_fallback = true
  enabled               = true
  exclude_office_ips    = false
  match                 = "identity.email == \"foo@example.com\""
  name                  = "xocmddmeyz"
  precedence            = 5
  support_url           = "support_url"
  switch_locked         = true
  service_mode_v2 = {
    mode = "warp"
  }
}

resource "cloudflare_zero_trust_device_custom_profile" "terraform_managed_resource_1" {
  account_id            = "f037e56e89293a057740de681ac9abbe"
  allow_mode_switch     = true
  allow_updates         = true
  allowed_to_leave      = true
  auto_connect          = 0
  captive_portal        = 5
  description           = "xeqtpkxdkw"
  disable_auto_fallback = true
  enabled               = true
  exclude_office_ips    = false
  match                 = "identity.email == \"foo@example.com\""
  name                  = "xeqtpkxdkw"
  precedence            = 10
  support_url           = "support_url"
  switch_locked         = true
  service_mode_v2 = {
    mode = "warp"
  }
}

