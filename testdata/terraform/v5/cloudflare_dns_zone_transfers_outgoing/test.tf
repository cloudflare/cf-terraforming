resource "cloudflare_dns_zone_transfers_outgoing" "terraform_managed_resource" {
  name    = "terraform.cfapi.net."
  peers   = ["226c39f30555498082ba9c4a3204e1b8"]
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
}

