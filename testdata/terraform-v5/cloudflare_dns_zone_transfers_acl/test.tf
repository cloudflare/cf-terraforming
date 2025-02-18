resource "cloudflare_dns_zone_transfers_acl" "example_dns_zone_transfers_acl" {
  account_id = "01a7362d577a6c3019a474fd6f485823"
  ip_range = "192.0.2.53/28"
  name = "my-acl-1"
}
