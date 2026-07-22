resource "cloudflare_regional_hostname" "terraform_managed_resource" {
  hostname   = "foo.example.com.terraform.cfapi.net"
  region_key = "ca"
  zone_id    = "0da42c8d2132a9ddaf714f9e7c920711"
}

