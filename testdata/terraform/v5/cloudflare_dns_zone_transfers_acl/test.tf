resource "cloudflare_dns_zone_transfers_acl" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  ip_range   = "1.2.3.4/32"
  name       = "jsjzqpgumk"
}

