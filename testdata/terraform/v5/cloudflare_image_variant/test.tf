resource "cloudflare_image_variant" "thumbnail" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  id         = "thumbnail"
  options    = {
    fit      = "scale-down"
    metadata = "keep"
    height   = 200
    width    = 200
  }
  never_require_signed_urls = false
}

