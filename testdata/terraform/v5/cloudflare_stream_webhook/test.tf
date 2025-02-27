resource "cloudflare_stream_webhook" "terraform_managed_resource" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  notification_url = "https://example.com"
}
