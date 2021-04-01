resource "cloudflare_load_balancer_pool" "terraform_managed_resource" {
  check_regions = [ "WEU", "ENAM" ]
  created_on = "2014-01-01T05:20:00.12345Z"
  description = "Primary data center - Provider XYZ"
  enabled = false
  minimum_origins = 2
  modified_on = "2014-01-01T05:20:00.12345Z"
  monitor = "f1aba936b94213e5b8dca0c0dbf1f9cc"
  name = "primary-dc-1"
  notification_email = "someone@example.com,sometwo@example.com"
}
