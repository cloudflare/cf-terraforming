resource "cloudflare_healthcheck" "terraform_managed_resource" {
  address               = "example.com"
  check_regions         = ["WNAM"]
  consecutive_fails     = 1
  consecutive_successes = 1
  interval              = 60
  name                  = "zngpvvwgvw"
  retries               = 2
  suspended             = false
  timeout               = 5
  type                  = "TCP"
  zone_id               = "0da42c8d2132a9ddaf714f9e7c920711"
  tcp_config = {
    method = "connection_established"
    port   = 80
  }
}

