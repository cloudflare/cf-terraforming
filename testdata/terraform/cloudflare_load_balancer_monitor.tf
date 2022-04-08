resource "cloudflare_load_balancer_monitor" "terraform_managed_resource" {
  allow_insecure   = true
  description      = "Login page monitor"
  expected_body    = "alive"
  expected_codes   = "2xx"
  follow_redirects = true
  interval         = 90
  method           = "GET"
  path             = "/health"
  port             = 8080
  probe_zone       = "example.com"
  retries          = 0
  timeout          = 3
  type             = "https"
  header {
    header = ""
    values = ""
  }
}
