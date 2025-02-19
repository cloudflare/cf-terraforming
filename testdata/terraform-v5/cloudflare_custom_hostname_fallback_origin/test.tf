resource "cloudflare_custom_hostname_fallback_origin" "terraform_managed_resource" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  origin = "fallback.example.com"
}
