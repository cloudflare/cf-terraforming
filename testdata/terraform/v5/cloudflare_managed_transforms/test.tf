resource "cloudflare_managed_transforms" "terraform_managed_resource" {
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  managed_request_headers = [{
    enabled = false
    id      = "add_bot_protection_headers"
  }]
  managed_response_headers = [{
    enabled = false
    id      = "remove_x-powered-by_header"
  }]
}

