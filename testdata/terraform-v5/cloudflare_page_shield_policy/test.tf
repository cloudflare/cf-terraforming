resource "cloudflare_page_shield_policy" "example_page_shield_policy" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  action = "allow"
  description = "Checkout page CSP policy"
  enabled = true
  expression = "ends_with(http.request.uri.path, \"/checkout\")"
  value = "script-src \'none\';"
}
