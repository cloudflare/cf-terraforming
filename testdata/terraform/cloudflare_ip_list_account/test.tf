resource "cloudflare_ip_list" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  name = "example_list"
  kind = "ip"
  description = "list description"

  item {
    value = "192.0.2.1"
    comment = "Office IP"
  }

  item {
    value = "203.0.113.0/24"
    comment = "Datacenter range"
  }
}
