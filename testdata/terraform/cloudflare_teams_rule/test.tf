resource "cloudflare_teams_rule" "terraform_managed_resource" {
  account_id     = "f037e56e89293a057740de681ac9abbe"
  action         = "block"
  description    = "Block bad websites based on their host name."
  device_posture = "any(device_posture.checks.passed[*] in {\"1308749e-fcfb-4ebc-b051-fe022b632644\"})"
  enabled        = true
  filters        = ["http"]
  identity       = "any(identity.groups.name[*] in {\"finance\"})"
  name           = "block bad websites"
  precedence     = 0
  traffic        = "http.request.uri matches \".*a/partial/uri.*\" and http.request.host in $01302951-49f9-47c9-a400-0297e60b6a10"
  rule_settings {
    add_headers        = {}
    allow_child_bypass = false
    audit_ssh {
      command_logging = false
    }
    block_page_enabled = true
    bypass_parent_rule = false
    check_session {
      duration = "5m0s"
      enforce  = true
    }
    egress {
      ipv4          = "192.0.2.2"
      ipv4_fallback = "192.0.2.3"
      ipv6          = "2001:DB8::/64"
    }
    insecure_disable_dnssec_validation = false
    ip_categories                      = true
    l4override {
      ip   = "1.1.1.1"
      port = 53
    }
    override_host = "example.com"
    override_ips  = ["1.1.1.1", "2.2.2.2"]
    payload_log {
      enabled = true
    }
    untrusted_cert {
      action = "error"
    }
  }
}