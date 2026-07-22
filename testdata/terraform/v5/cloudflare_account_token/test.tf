resource "cloudflare_account_token" "vpd-tests" {
  account_id = "a67e14daa5f8dceeb91fe5449ba496eb"
  name = "readonly token"
  policies = [
    {
      effect = "allow"
      resources = "*"
      permission_groups = [{
        id = "e53155bef40b42b1b150c9d5700e29fa",
      }]
    },
  ]
}
