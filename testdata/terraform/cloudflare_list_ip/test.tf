resource "cloudflare_list" "terraform_managed_resource" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  description = "This is a note"
  kind        = "ip"
  name        = "ip_list"
  item {
    comment = "one"
    value {
      ip = "10.0.0.1"
    }
  }
  item {
    comment = "two"
    value {
      ip = "10.0.0.2"
    }
  }
}
