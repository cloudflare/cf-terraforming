resource "cloudflare_regional_tiered_cache" "terraform_managed_resource" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  value = "on"
}
