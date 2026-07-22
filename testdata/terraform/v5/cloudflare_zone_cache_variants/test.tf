resource "cloudflare_zone_cache_variants" "terraform_managed_resource" {
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  value = {
    avif = ["image/avif", "image/webp"]
    bmp  = ["image/bmp", "image/webp"]
    gif  = ["image/gif", "image/webp"]
    jp2  = ["image/jp2", "image/webp"]
    jpeg = ["image/jpeg", "image/webp"]
    jpg  = ["image/jpg", "image/webp"]
    jpg2 = ["image/jpg2", "image/webp"]
    png  = ["image/png"]
    tif  = ["image/tif", "image/webp"]
    tiff = ["image/tiff", "image/webp"]
    webp = ["image/webp"]
  }
}

