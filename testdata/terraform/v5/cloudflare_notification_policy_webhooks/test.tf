resource "cloudflare_notification_policy_webhooks" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  name       = "my webhooks destination for receiving Cloudflare notifications"
  url        = "https://httpbin.org/post"
}

