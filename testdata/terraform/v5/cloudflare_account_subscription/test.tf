resource "cloudflare_account_subscription" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  frequency  = "monthly"
  rate_plan = {
    currency           = "USD"
    externally_managed = false
    id                 = "enterprise"
    is_contract        = true
    public_name        = "Image Resizing Ent"
    scope              = "user"
    sets               = ["usage", "is_cloudflare", "public"]
  }
}