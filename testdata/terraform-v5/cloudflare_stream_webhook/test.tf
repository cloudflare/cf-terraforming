resource "cloudflare_stream_webhook" "example_stream_webhook" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  notification_url = "https://example.com"
}
