resource "cloudflare_zero_trust_access_application" "terraform_managed_resource" {
  account_id                 = "f037e56e89293a057740de681ac9abbe"
  allowed_idps               = []
  app_launcher_visible       = true
  auto_redirect_to_identity  = false
  domain                     = "gpfqbfyfcx.terraform.cfapi.net"
  enable_binding_cookie      = false
  http_only_cookie_attribute = true
  name                       = "gpfqbfyfcx"
  options_preflight_bypass   = false
  session_duration           = "24h"
  tags                       = []
  type                       = "self_hosted"
  destinations = [{
    type = "public"
    uri  = "gpfqbfyfcx.terraform.cfapi.net"
  }]
  policies = []
}

