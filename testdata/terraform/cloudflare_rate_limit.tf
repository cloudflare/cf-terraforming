resource "cloudflare_rate_limit" "terraform_managed_resource" {
  bypass_url_patterns {
    name = "url"
    value = "example.com/allowed-bypass"
  }
  description = "example rate limit"
  period = 60
  threshold = 10
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  action {
    response {
      body = "{\"response\":\"your request has been rate limited\"}"
      content_type = "application/json"
    }
    mode = "ban"
    timeout = 3600
  }
  match {
    request {
      methods = [ "POST" ]
      schemes = [ "_ALL_" ]
      url_pattern = "example.com"
    }
    response {
      headers {
        name = "My_origin_field"
        op = "eq"
        value = "block_request"
      }
      origin_traffic = false
      statuses = [ 401, 403 ]
    }
  }
}
