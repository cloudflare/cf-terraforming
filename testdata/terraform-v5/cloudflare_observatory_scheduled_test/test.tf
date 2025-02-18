resource "cloudflare_observatory_scheduled_test" "example_observatory_scheduled_test" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  url = "example.com"
}
