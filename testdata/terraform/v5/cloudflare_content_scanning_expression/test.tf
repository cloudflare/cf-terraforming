resource "cloudflare_content_scanning_expression" "terraform_managed_resource" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  body = [{
    payload = "lookup_json_string(http.request.body.raw, \"file\")"
  }]
}
