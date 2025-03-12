resource "cloudflare_web3_hostname" "terraform_managed_resource" {
  description = "test"
  name        = "mickmslsyi.terraform.cfapi.net"
  target      = "ethereum"
  zone_id     = "0da42c8d2132a9ddaf714f9e7c920711"
}

