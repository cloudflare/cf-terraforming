resource "cloudflare_access_identity_provider" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  name = "PIN login"
  type = "onetimepin"
}
