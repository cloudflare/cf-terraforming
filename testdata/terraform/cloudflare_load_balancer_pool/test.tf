resource "cloudflare_load_balancer_pool" "terraform_managed_resource" {
  account_id         = "f037e56e89293a057740de681ac9abbe"
  check_regions      = ["WEU", "ENAM"]
  description        = "Primary data center - Provider XYZ"
  enabled            = false
  minimum_origins    = 2
  monitor            = "f1aba936b94213e5b8dca0c0dbf1f9cc"
  name               = "primary-dc-1"
  notification_email = "someone@example.com,sometwo@example.com"
  origins {
    address = "0.0.0.0"
    enabled = true
    name    = "app-server-1"
    weight  = 1
    header {
      header = "Host"
      values = ["example.com"]
    }
  }
}
