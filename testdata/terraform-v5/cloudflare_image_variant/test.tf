resource "cloudflare_image_variant" "example_image_variant" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  id = "hero"
  options = {
    fit = "scale-down"
    height = 768
    metadata = "keep"
    width = 1366
  }
  never_require_signed_urls = true
}
