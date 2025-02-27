resource "cloudflare_load_balancer_monitor" "terraform_managed_resource" {
  account_id       = "f037e56e89293a057740de681ac9abbe"
  allow_insecure   = true
  consecutive_down = 2
  consecutive_up   = 3
  description      = "Login page monitor"
  expected_body    = "alive"
  expected_codes   = "2xx"
  follow_redirects = true
  interval         = 90
  method           = "GET"
  path             = "/health"
  port             = 8080
  probe_zone       = "example.com"
  retries          = 1
  timeout          = 3
  type             = "https"
}
