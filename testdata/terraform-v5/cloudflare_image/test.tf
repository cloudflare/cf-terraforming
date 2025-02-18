resource "cloudflare_image" "example_image" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  file = {

  }
  metadata = {

  }
  require_signed_urls = true
  url = "https://example.com/path/to/logo.png"
}
