resource "cloudflare_api_shield" "terraform_managed_resource" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  auth_id_characteristics = [{
    name = "authorization"
    type = "header"
  }]
}
