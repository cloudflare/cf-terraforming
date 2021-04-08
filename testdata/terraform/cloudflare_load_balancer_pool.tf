resource "cloudflare_load_balancer_pool" "terraform_managed_resource" {
  check_regions = [ "WEU", "ENAM" ]
  description = "Primary data center - Provider XYZ"
  enabled = false
  minimum_origins = 2
  monitor = "f1aba936b94213e5b8dca0c0dbf1f9cc"
  name = "primary-dc-1"
  notification_email = "someone@example.com,sometwo@example.com"
}
