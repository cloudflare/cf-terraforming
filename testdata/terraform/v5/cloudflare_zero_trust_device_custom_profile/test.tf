resource "cloudflare_zero_trust_device_custom_profile" "terraform_managed_resource_0" {
  account_id                     = "f037e56e89293a057740de681ac9abbe"
  allow_mode_switch              = true
  allow_updates                  = true
  allowed_to_leave               = true
  auto_connect                   = 0
  captive_portal                 = 5
  description                    = "xocmddmeyz"
  disable_auto_fallback          = true
  enabled                        = true
  exclude_office_ips             = false
  match                          = "identity.email == \"foo@example.com\""
  name                           = "xocmddmeyz"
  precedence                     = 5
  register_interface_ip_with_dns = true
  support_url                    = "support_url"
  switch_locked                  = true
  exclude = [{
    address = "10.0.0.0/8"
    }, {
    address = "100.64.0.0/10"
    }, {
    address     = "169.254.0.0/16"
    description = "DHCP Unspecified"
    }, {
    address = "172.16.0.0/12"
    }, {
    address = "192.0.0.0/24"
    }, {
    address = "192.168.0.0/16"
    }, {
    address = "224.0.0.0/24"
    }, {
    address = "240.0.0.0/4"
    }, {
    address     = "255.255.255.255/32"
    description = "DHCP Broadcast"
    }, {
    address     = "fe80::/10"
    description = "IPv6 Link Local"
    }, {
    address = "fd00::/8"
    }, {
    address = "ff01::/16"
    }, {
    address = "ff02::/16"
    }, {
    address = "ff03::/16"
    }, {
    address = "ff04::/16"
    }, {
    address = "ff05::/16"
  }]
  service_mode_v2 = {
    mode = "warp"
  }
}

resource "cloudflare_zero_trust_device_custom_profile" "terraform_managed_resource_1" {
  account_id                     = "f037e56e89293a057740de681ac9abbe"
  allow_mode_switch              = true
  allow_updates                  = true
  allowed_to_leave               = true
  auto_connect                   = 0
  captive_portal                 = 5
  description                    = "xeqtpkxdkw"
  disable_auto_fallback          = true
  enabled                        = true
  exclude_office_ips             = false
  match                          = "identity.email == \"foo@example.com\""
  name                           = "xeqtpkxdkw"
  precedence                     = 10
  register_interface_ip_with_dns = true
  support_url                    = "support_url"
  switch_locked                  = true
  exclude = [{
    address = "10.0.0.0/8"
    }, {
    address = "100.64.0.0/10"
    }, {
    address     = "169.254.0.0/16"
    description = "DHCP Unspecified"
    }, {
    address = "172.16.0.0/12"
    }, {
    address = "192.0.0.0/24"
    }, {
    address = "192.168.0.0/16"
    }, {
    address = "224.0.0.0/24"
    }, {
    address = "240.0.0.0/4"
    }, {
    address     = "255.255.255.255/32"
    description = "DHCP Broadcast"
    }, {
    address     = "fe80::/10"
    description = "IPv6 Link Local"
    }, {
    address = "fd00::/8"
    }, {
    address = "ff01::/16"
    }, {
    address = "ff02::/16"
    }, {
    address = "ff03::/16"
    }, {
    address = "ff04::/16"
    }, {
    address = "ff05::/16"
  }]
  service_mode_v2 = {
    mode = "warp"
  }
}

