resource "cloudflare_access_application" "terraform_managed_resource" {
  account_id       = "f037e56e89293a057740de681ac9abbe"
  allowed_idps     = ["699d98642c564d2e855e9661899b7252"]
  domain           = "test.example.com/admin"
  name             = "Admin Site"
  session_duration = "24h"
  foo              = "bar"
}
