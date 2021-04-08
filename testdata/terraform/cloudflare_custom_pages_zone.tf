resource "cloudflare_custom_pages" "terraform_managed_resource" {
  state = "customized"
  type = "basic_challenge"
  url = "https://example.com/challenge.html"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
}
