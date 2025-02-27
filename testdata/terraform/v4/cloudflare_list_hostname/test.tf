resource "cloudflare_list" "terraform_managed_resource" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  description = "This is a note"
  kind        = "hostname"
  name        = "hostname_list"
  item {
    comment = "one"
    value {
      hostname {
        url_hostname = "example.com"
      }
    }
  }
}
