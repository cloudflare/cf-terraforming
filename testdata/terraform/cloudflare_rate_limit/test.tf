resource "cloudflare_rate_limit" "terraform_managed_resource" {
  bypass_url_patterns = ["example.com/allowed-bypass", "example.com/allowed-bypass-other"]
  description         = "example rate limit"
  period              = 60
  threshold           = 10
  zone_id             = "0da42c8d2132a9ddaf714f9e7c920711"
  action {
    mode    = "ban"
    timeout = 3600
    response {
      body         = "{\"response\":\"your request has been rate limited\"}"
      content_type = "application/json"
    }
  }
  match {
    request {
      methods     = ["POST"]
      schemes     = ["_ALL_"]
      url_pattern = "example.com"
    }
    response {
      headers = [
        {
          name  = "My_origin_field"
          op    = "eq"
          value = "block_request"
        },
        {
          name  = "Other"
          op    = "eq"
          value = "block_request"
        }
      ]
      origin_traffic = false
      statuses       = [401, 403]
    }
  }
}
