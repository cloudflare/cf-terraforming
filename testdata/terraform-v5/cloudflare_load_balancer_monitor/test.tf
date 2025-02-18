resource "cloudflare_load_balancer_monitor" "example_load_balancer_monitor" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  allow_insecure = true
  consecutive_down = 0
  consecutive_up = 0
  description = "Login page monitor"
  expected_body = "alive"
  expected_codes = "2xx"
  follow_redirects = true
  header = {
    Host = ["example.com"]
    X-App-ID = ["abc123"]
  }
  interval = 0
  method = "GET"
  path = "/health"
  port = 0
  probe_zone = "example.com"
  retries = 0
  timeout = 0
  type = "http"
}
