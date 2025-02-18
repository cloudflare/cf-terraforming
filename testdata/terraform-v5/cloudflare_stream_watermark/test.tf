resource "cloudflare_stream_watermark" "example_stream_watermark" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  file = "@/Users/rchen/Downloads/watermark.png"
  name = "Marketing Videos"
  opacity = 0.75
  padding = 0.1
  position = "center"
  scale = 0.1
}
