resource "cloudflare_dns_zone_transfers_peer" "example_dns_zone_transfers_peer" {
  account_id = "01a7362d577a6c3019a474fd6f485823"
  name = "my-peer-1"
}
