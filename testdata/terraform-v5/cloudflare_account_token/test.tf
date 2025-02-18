resource "cloudflare_account_token" "example_account_token" {
  account_id = "eb78d65290b24279ba6f44721b3ea3c4"
  name = "readonly token"
  policies = [{
    effect = "allow"
    permission_groups = [{
      meta = {
        key = "key"
        value = "value"
      }
    }, {
      meta = {
        key = "key"
        value = "value"
      }
    }]
    resources = {
      "com.cloudflare.api.account.zone.22b1de5f1c0e4b3ea97bb1e963b06a43" = "*"
      "com.cloudflare.api.account.zone.eb78d65290b24279ba6f44721b3ea3c4" = "*"
    }
  }]
  condition = {
    request_ip = {
      in = ["123.123.123.0/24", "2606:4700::/32"]
      not_in = ["123.123.123.100/24", "2606:4700:4700::/48"]
    }
  }
  expires_on = "2020-01-01T00:00:00Z"
  not_before = "2018-07-01T05:20:00Z"
}
