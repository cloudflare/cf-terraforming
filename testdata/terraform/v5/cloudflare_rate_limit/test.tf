resource "cloudflare_rate_limit" "terraform_managed_resource" {
  period    = 900
  threshold = 60
  zone_id   = "0da42c8d2132a9ddaf714f9e7c920711"
  action = {
    mode = "ban"
    response = {
      body         = "{\"response\":\"your request has been rate limited\"}"
      content_type = "application/json"
    }
    timeout = 3600
  }
  match = {
    request = {
      methods = ["POST"]
      schemes = ["_ALL_"]
      url     = "example.com"
    }
    response = {
      headers = [{
        name  = "My_origin_field"
        op    = "eq"
        value = "block_request"
        }, {
        name  = "Other"
        op    = "eq"
        value = "block_request"
      }]
      origin_traffic = false
      status         = [401, 403]
    }
  }
}

