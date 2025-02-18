resource "cloudflare_web_analytics_rule" "example_web_analytics_rule" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  ruleset_id = "f174e90a-fafe-4643-bbbc-4a0ed4fc8415"
  host = "example.com"
  inclusive = true
  is_paused = false
  paths = ["*"]
}
