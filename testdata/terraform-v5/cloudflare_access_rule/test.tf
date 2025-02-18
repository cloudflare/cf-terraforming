resource "cloudflare_access_rule" "example_access_rule" {
  configuration = {
    target = "ip"
    value = "198.51.100.4"
  }
  mode = "block"
  zone_id = "zone_id"
  notes = "This rule is enabled because of an event that occurred on date X."
}
