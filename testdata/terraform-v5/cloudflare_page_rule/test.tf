resource "cloudflare_page_rule" "example_page_rule" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  priority = 0
  status = "active"
}
