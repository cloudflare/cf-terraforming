resource "cloudflare_api_shield" "terraform_managed_resource" {
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  auth_id_characteristics {
    name = "test-header"
    type = "header"
  }
  auth_id_characteristics {
    name = "test-cookie"
    type = "cookie"
  }
}
