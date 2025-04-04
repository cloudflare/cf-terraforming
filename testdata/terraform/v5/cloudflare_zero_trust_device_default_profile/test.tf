resource "cloudflare_zero_trust_device_default_profile" "terraform_managed_resource" {
  account_id            = "f037e56e89293a057740de681ac9abbe"
  allow_mode_switch     = true
  allow_updates         = true
  allowed_to_leave      = true
  auto_connect          = 0
  captive_portal        = 5
  disable_auto_fallback = true
  exclude_office_ips    = true
  support_url           = "https://cloudflare.com"
  switch_locked         = true
  tunnel_protocol       = "wireguard"
  service_mode_v2 = {
    mode = "warp"
  }
}

