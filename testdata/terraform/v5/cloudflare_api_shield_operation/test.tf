resource "cloudflare_api_shield_operation" "terraform_managed_resource" {
  endpoint = "/example/path"
  host     = "terraform.cfapi.net"
  method   = "GET"
  zone_id  = "0da42c8d2132a9ddaf714f9e7c920711"
}
