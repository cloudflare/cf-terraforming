resource "cloudflare_observatory_scheduled_test" "terraform_managed_resource" {
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  url     = urlencode("terraform.cfapi.net/thyygxveip")
}

