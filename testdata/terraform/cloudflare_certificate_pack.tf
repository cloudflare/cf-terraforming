resource "cloudflare_certificate_pack" "terraform_managed_resource" {
  hosts = [ "example.com", "*.example.com", "www.example.com" ]
  type = "dedicated_custom"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
}
