resource "cloudflare_load_balancer_monitor" "terraform_managed_resource" {
  allow_insecure = true
  created_on = "2014-01-01T05:20:00.12345Z"
  description = "Login page monitor"
  expected_body = "alive"
  expected_codes = "2xx"
  follow_redirects = true
  interval = 90
  method = "GET"
  modified_on = "2014-01-01T05:20:00.12345Z"
  path = "/health"
  port = 8080
  probe_zone = "example.com"
  retries = 0
  timeout = 3
  type = "https"
}
