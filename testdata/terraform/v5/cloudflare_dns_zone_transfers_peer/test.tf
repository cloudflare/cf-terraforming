resource "cloudflare_dns_zone_transfers_peer" "terraform_managed_resource" {
  account_id = "01a7362d577a6c3019a474fd6f485823"
  name = "my-peer-1"
}
