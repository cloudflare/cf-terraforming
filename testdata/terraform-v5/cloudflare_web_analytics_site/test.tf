resource "cloudflare_web_analytics_site" "example_web_analytics_site" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  auto_install = true
  host = "example.com"
  zone_tag = "023e105f4ecef8ad9ca31a8372d0c353"
}
