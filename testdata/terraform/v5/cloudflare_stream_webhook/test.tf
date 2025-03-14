resource "cloudflare_stream_webhook" "terraform_managed_resource" {
  account_id       = "f037e56e89293a057740de681ac9abbe"
  notification_url = "https://example.com"
}

