resource "cloudflare_rate_limit" "example_rate_limit" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  action = {
    mode = "simulate"
    response = {
      body = "<error>This request has been rate-limited.</error>"
      content_type = "text/xml"
    }
    timeout = 86400
  }
  match = {
    headers = [{
      name = "Cf-Cache-Status"
      op = "eq"
      value = "HIT"
    }]
    request = {
      methods = ["GET", "POST"]
      schemes = ["HTTP", "HTTPS"]
      url = "*.example.org/path*"
    }
    response = {
      origin_traffic = true
    }
  }
  period = 900
  threshold = 60
}
