resource "cloudflare_notification_policy_webhooks" "terraform_managed_resource" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  name = "Slack Webhook"
  url = "https://hooks.slack.com/services/Ds3fdBFbV/456464Gdd"
  secret = "secret"
}
