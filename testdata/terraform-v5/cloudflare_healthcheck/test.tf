resource "cloudflare_healthcheck" "example_healthcheck" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  address = "www.example.com"
  name = "server-1"
  check_regions = ["WNAM", "ENAM"]
  consecutive_fails = 0
  consecutive_successes = 0
  description = "Health check for www.example.com"
  http_config = {
    allow_insecure = true
    expected_body = "success"
    expected_codes = ["2xx", "302"]
    follow_redirects = true
    header = {
      Host = ["example.com"]
      X-App-ID = ["abc123"]
    }
    method = "GET"
    path = "/health"
    port = 0
  }
  interval = 0
  retries = 0
  suspended = true
  tcp_config = {
    method = "connection_established"
    port = 0
  }
  timeout = 0
  type = "HTTPS"
}
