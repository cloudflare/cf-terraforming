resource "cloudflare_custom_hostname_fallback_origin" "terraform_managed_resource" {
  origin  = "fallback-origin.jtwuheolsi.terraform.cfapi.net"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
}
