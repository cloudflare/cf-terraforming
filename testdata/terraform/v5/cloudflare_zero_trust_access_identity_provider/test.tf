resource "cloudflare_zero_trust_access_identity_provider" "terraform_managed_resource_0" {
  name    = "Widget Corps IDP"
  type    = "onetimepin"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  config  = {}
  scim_config = {
    enabled                  = false
    group_member_deprovision = false
    seat_deprovision         = false
    user_deprovision         = false
  }
}

resource "cloudflare_zero_trust_access_identity_provider" "terraform_managed_resource_1" {
  name    = "lnfbpxpksi"
  type    = "azureAD"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  config = {
    azure_cloud                = "default"
    client_id                  = "test"
    conditional_access_enabled = false
    directory_id               = "directory"
    support_groups             = true
  }
  scim_config = {
    enabled                  = true
    group_member_deprovision = false
    seat_deprovision         = true
    user_deprovision         = true
  }
}

