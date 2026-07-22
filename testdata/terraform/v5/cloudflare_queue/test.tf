resource "cloudflare_queue" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  queue_name = "test-q"
  settings = {
    delivery_delay           = 0
    message_retention_period = 345600
  }
}

