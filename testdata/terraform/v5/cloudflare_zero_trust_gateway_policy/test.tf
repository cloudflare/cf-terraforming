resource "cloudflare_zero_trust_gateway_policy" "terraform_managed_resource_0" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  action      = "block"
  description = "desc"
  enabled     = false
  filters     = ["dns"]
  name        = "rytdytpfmz"
  precedence  = 12302
  traffic     = "any(dns.domains[*] == \"example.com\")"
  rule_settings = {
    add_headers         = null
    biso_admin_controls = null
    block_page_enabled  = true
    block_reason        = "cuz"
    check_session       = null
    egress = {
      ipv4 = "203.0.113.1"
      ipv6 = "2001:db8::/32"
    }
    insecure_disable_dnssec_validation = false
    ip_categories                      = false
    ip_indicator_feeds                 = false
    l4override                         = null
    override_host                      = ""
    override_ips                       = null
    payload_log = {
      enabled = true
    }
    untrusted_cert = {
      action = "error"
    }
  }
}

resource "cloudflare_zero_trust_gateway_policy" "terraform_managed_resource_1" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  action      = "block"
  description = "test description "
  enabled     = true
  filters     = ["http"]
  name        = "test-http"
  precedence  = 13302
  traffic     = "http.request.uri.path == \"/foo\""
  rule_settings = {
    add_headers                        = null
    biso_admin_controls                = null
    block_page_enabled                 = false
    block_reason                       = ""
    check_session                      = null
    insecure_disable_dnssec_validation = false
    ip_categories                      = false
    ip_indicator_feeds                 = false
    l4override                         = null
    override_host                      = ""
    override_ips                       = null
  }
}

