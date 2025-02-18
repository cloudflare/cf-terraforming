resource "cloudflare_tiered_cache" "example_tiered_cache" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  value = "on"
}
