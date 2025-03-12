resource "cloudflare_dns_zone_transfers_incoming" "terraform_managed_resource" {
  auto_refresh_seconds = 300
  name                 = "terraform.cfapi.net."
  peers                = ["e77bdc034b754e2fbdc622fed1cf6b92"]
  zone_id              = "0da42c8d2132a9ddaf714f9e7c920711"
}

