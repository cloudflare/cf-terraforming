resource "cloudflare_zone_cache_variants" "example_zone_cache_variants" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  value = {
    avif = ["image/webp", "image/jpeg"]
    bmp = ["image/webp", "image/jpeg"]
    gif = ["image/webp", "image/jpeg"]
    jp2 = ["image/webp", "image/avif"]
    jpeg = ["image/webp", "image/avif"]
    jpg = ["image/webp", "image/avif"]
    jpg2 = ["image/webp", "image/avif"]
    png = ["image/webp", "image/avif"]
    tif = ["image/webp", "image/avif"]
    tiff = ["image/webp", "image/avif"]
    webp = ["image/jpeg", "image/avif"]
  }
}
