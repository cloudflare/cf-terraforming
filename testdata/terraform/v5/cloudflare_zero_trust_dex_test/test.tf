resource "cloudflare_zero_trust_dex_test" "terraform_managed_resource_0" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  description = "qmktnyxzho"
  enabled     = true
  interval    = "0h30m0s"
  name        = "qmktnyxzho"
  data = {
    host = "1.1.1.1"
    kind = "traceroute"
  }
  target_policies = []
}

resource "cloudflare_zero_trust_dex_test" "terraform_managed_resource_1" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  description = "aibzxpxpyl"
  enabled     = true
  interval    = "0h30m0s"
  name        = "aibzxpxpyl"
  data = {
    host = "foo.cloudflare.com"
    kind = "traceroute"
  }
  target_policies = []
}

