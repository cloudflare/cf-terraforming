resource "cloudflare_authenticated_origin_pulls" "terraform_managed_resource" {
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  config = [{
    cert_id  = "0a96490d-0bec-4ef6-b701-99f19f28d320"
    enabled  = false
    hostname = "jotsqcjaho.terraform.cfapi.net"
  }]
}

