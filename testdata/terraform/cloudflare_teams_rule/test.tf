resource "cloudflare_teams_rule" "terraform_managed_resource" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  name        = "block bad websites"
  description = "Block bad websites based on their host name."
  precedence  = 0
  action      = "block"
  filters     = ["http"]
  device_posture = "any(device_posture.checks.passed[*] in {\"1308749e-fcfb-4ebc-b051-fe022b632644\"})"
  identity = "any(identity.groups.name[*] in {\"finance\"})"
  traffic     = "http.request.uri matches \".*a/partial/uri.*\" and http.request.host in $01302951-49f9-47c9-a400-0297e60b6a10"
  rule_settings  {
    add_headers = {
      "x-Custom-Header-Name" : "somecustomvalue"
    }
    allow_child_bypass = false
    audit_ssh {
      command_logging = false
    }
    biso_admin_controls {
      disable_copy_paste = false
      disable_download = false
      disable_keyboard = false
      disable_printing = false
      disable_upload = false
    }
    block_page_enabled = true
    block_page_reason  = "This website is a security risk"
    bypass_parent_rule = true
    check_session {
      duration ="300s"
      enforce = true
    }
    egress {
      ipv4 = "192.0.2.2"
      ipv6 = "2001:DB8::/64"
      ipv4_fallback = "192.0.2.3"
    }
    insecure_disable_dnssec_validation = false
    ip_categories = true
    l4override {
      ip = "1.1.1.1"
      port = 0
    }
    override_host = "example.com"
    override_ips = ["1.1.1.1", "2.2.2.2"]
    payload_log {
      enabled = true
    }
    untrusted_cert {
      action = "error"
    }
  }
}