resource "cloudflare_zero_trust_access_group" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  name       = "kskryelbix"
  exclude    = []
  include = [{
    email_list = {
      id = "140f88f4-e69f-4d51-8ca0-8af660000039"
    }
  }]
  require = []
}

