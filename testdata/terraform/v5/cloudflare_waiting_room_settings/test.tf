resource "cloudflare_waiting_room_settings" "terraform_managed_resource" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  search_engine_crawler_bypass = true
}
