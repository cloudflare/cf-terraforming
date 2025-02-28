resource "cloudflare_account_member" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  roles = [{
    description = "Can edit any Cloudflare setting, make purchases, update billing, and manage memberships. Super Administrators can revoke the access of other Super Administrators."
    id          = "33666b9c79b9a5273fc7344ff42f953d"
    name        = "Super Administrator - All Privileges"
    permissions = {
      access = {
        edit = true
        read = true
      },
      analytics = {
        edit = false
        read = true
      }
    }
  }]
  status = "accepted"
  policies = [{
    access            = "allow"
    id                = "aa0fe1963f8f49c8a4f2d73314146e23"
    permission_groups = {
      id = "68dfcdfb261c4d48aa75ec0e5413f2c3"
      meta = {
        description = "Can edit any Cloudflare setting, make purchases, update billing, and manage memberships. Super Administrators can revoke the access of other Super Administrators."
        editable    = false
        label       = "all_privileges"
        scopes      = "com.cloudflare.api.account"
        name        = "Super Administrator - All Privileges"
      }
    }
    resource_groups   = {
      id = "68dfcdfb261c4d48aa75ec0e5413f2c3"
      meta = {
        editable    = false
        name = "com.cloudflare.api.account.f037e56e89293a057740de681ac9abbe"
        scope = {
          key = "com.cloudflare.api.account.f037e56e89293a057740de681ac9abbe"
          objects = {
            key = "*"
          }
        }
      }
    }
  }]
}
