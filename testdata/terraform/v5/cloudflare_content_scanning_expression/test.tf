resource "cloudflare_content_scanning_expression" "terraform_managed_resource_0" {
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  body = [{
    payload = "lookup_json_string(http.request.body.raw, \"file\")"
  }]
}

resource "cloudflare_content_scanning_expression" "terraform_managed_resource_1" {
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  body = [{
    payload = "lookup_json_string(http.request.body.raw, \"file\")"
  }]
}

