resource "cloudflare_filter" "terraform_managed_resource" {
  expression = "(http.request.uri.path eq \"/hello\")"
  zone_id    = "0da42c8d2132a9ddaf714f9e7c920711"
}

