resource "cloudflare_zero_trust_device_posture_integration" "example_zero_trust_device_posture_integration" {
  account_id = "699d98642c564d2e855e9661899b7252"
  config = {
    api_url = "https://as123.awmdm.com/API"
    auth_url = "https://na.uemauth.vmwservices.com/connect/token"
    client_id = "example client id"
    client_secret = "example client secret"
  }
  interval = "10m"
  name = "My Workspace One Integration"
  type = "workspace_one"
}
