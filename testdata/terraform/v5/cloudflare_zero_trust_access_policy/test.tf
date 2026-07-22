resource "cloudflare_zero_trust_access_policy" "terraform_managed_resource_0" {
  account_id       = "f037e56e89293a057740de681ac9abbe"
  decision         = "non_identity"
  name             = "dryjkmkfpz"
  session_duration = "24h"
  exclude          = []
  include = [{
    ip = {
      ip = "127.0.0.1/32"
    }
  }]
  require = []
}

resource "cloudflare_zero_trust_access_policy" "terraform_managed_resource_1" {
  account_id       = "f037e56e89293a057740de681ac9abbe"
  decision         = "non_identity"
  name             = "dryjkmkfpz"
  session_duration = "24h"
  exclude          = []
  include = [{
    ip = {
      ip = "127.0.0.1/32"
    }
  }]
  require = []
}

