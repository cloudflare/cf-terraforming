resource "cloudflare_managed_headers" "terraform_managed_resource" {
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  managed_request_headers {
    enabled = true
    id      = "add_visitor_location_headers"
  }
  managed_response_headers {
    enabled = true
    id      = "remove_x-powered-by_header"
  }
}