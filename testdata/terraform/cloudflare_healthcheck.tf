resource "cloudflare_healthcheck" "terraform_managed_resource" {
  address = "www.example.com"
  check_regions = [ "WEU", "ENAM" ]
  consecutive_fails = 1
  consecutive_successes = 1
  description = "Health check for www.example.com"
  interval = 60
  name = "server-1"
  retries = 2
  suspended = false
  timeout = 5
  type = "HTTPS"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
}
