resource "cloudflare_zero_trust_gateway_policy" "example_zero_trust_gateway_policy" {
  account_id = "699d98642c564d2e855e9661899b7252"
  action = "on"
  name = "block bad websites"
  description = "Block bad websites based on their host name."
  device_posture = "any(device_posture.checks.passed[*] in {\"1308749e-fcfb-4ebc-b051-fe022b632644\"})"
  enabled = true
  expiration = {
    expires_at = "2014-01-01T05:20:20Z"
    duration = 10
    expired = false
  }
  filters = ["http"]
  identity = "any(identity.groups.name[*] in {\"finance\"})"
  precedence = 0
  rule_settings = {
    add_headers = {
      foo = "string"
    }
    allow_child_bypass = false
    audit_ssh = {
      command_logging = false
    }
    biso_admin_controls = {
      copy = "enabled"
      dcp = false
      dd = false
      dk = false
      download = "enabled"
      dp = false
      du = false
      keyboard = "enabled"
      paste = "enabled"
      printing = "enabled"
      upload = "enabled"
      version = "v1"
    }
    block_page_enabled = true
    block_reason = "This website is a security risk"
    bypass_parent_rule = false
    check_session = {
      duration = "300s"
      enforce = true
    }
    dns_resolvers = {
      ipv4 = [{
        ip = "2.2.2.2"
        port = 5053
        route_through_private_network = true
        vnet_id = "f174e90a-fafe-4643-bbbc-4a0ed4fc8415"
      }]
      ipv6 = [{
        ip = "2001:DB8::"
        port = 5053
        route_through_private_network = true
        vnet_id = "f174e90a-fafe-4643-bbbc-4a0ed4fc8415"
      }]
    }
    egress = {
      ipv4 = "192.0.2.2"
      ipv4_fallback = "192.0.2.3"
      ipv6 = "2001:DB8::/64"
    }
    ignore_cname_category_matches = true
    insecure_disable_dnssec_validation = false
    ip_categories = true
    ip_indicator_feeds = true
    l4override = {
      ip = "1.1.1.1"
      port = 0
    }
    notification_settings = {
      enabled = true
      msg = "msg"
      support_url = "support_url"
    }
    override_host = "example.com"
    override_ips = ["1.1.1.1", "2.2.2.2"]
    payload_log = {
      enabled = true
    }
    quarantine = {
      file_types = ["exe"]
    }
    resolve_dns_internally = {
      fallback = "none"
      view_id = "view_id"
    }
    resolve_dns_through_cloudflare = true
    untrusted_cert = {
      action = "pass_through"
    }
  }
  schedule = {
    fri = "08:00-12:30,13:30-17:00"
    mon = "08:00-12:30,13:30-17:00"
    sat = "08:00-12:30,13:30-17:00"
    sun = "08:00-12:30,13:30-17:00"
    thu = "08:00-12:30,13:30-17:00"
    time_zone = "America/New York"
    tue = "08:00-12:30,13:30-17:00"
    wed = "08:00-12:30,13:30-17:00"
  }
  traffic = "http.request.uri matches \".*a/partial/uri.*\" and http.request.host in $01302951-49f9-47c9-a400-0297e60b6a10"
}
