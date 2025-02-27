resource "cloudflare_access_application" "terraform_managed_resource" {
  account_id                = "f037e56e89293a057740de681ac9abbe"
  allowed_idps              = ["699d98642c564d2e855e9661899b7252"]
  auto_redirect_to_identity = false
  domain                    = "test.example.com/admin"
  enable_binding_cookie     = false
  name                      = "Admin Site"
  session_duration          = "24h"
}
