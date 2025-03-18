resource "cloudflare_zero_trust_device_posture_integration" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  interval   = "24h"
  name       = "ptstsjfmxl"
  type       = "workspace_one"
  config = {
    api_url   = "https://techp-as.awmdm.com/API"
    auth_url  = "https://na.uemauth.vmwservices.com/connect/token"
    client_id = "d0ed71f01c884e8b94ec4e4d6639f609"
  }
}

