resource "cloudflare_list" "terraform_managed_resource" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  description = "This is a note"
  kind        = "asn"
  name        = "asn_list"
  item {
    comment = "one"
    value {
      asn = 123
    }
  }
  item {
    comment = "two"
    value {
      asn = 456
    }
  }
}
