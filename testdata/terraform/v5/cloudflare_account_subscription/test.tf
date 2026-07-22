resource "cloudflare_account_subscription" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  frequency  = "monthly"
  rate_plan = {
    id    = "enterprise"
    scope = "user"
  }
}