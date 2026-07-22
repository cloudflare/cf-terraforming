resource "cloudflare_dns_zone_transfers_peer" "terraform_managed_resource_0" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  ip          = "1.2.3.4"
  ixfr_enable = false
  name        = "terraform-peer"
  port        = 53
}

resource "cloudflare_dns_zone_transfers_peer" "terraform_managed_resource_1" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  ixfr_enable = false
  name        = "terraform-peer"
  port        = 0
}

resource "cloudflare_dns_zone_transfers_peer" "terraform_managed_resource_2" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  ixfr_enable = false
  name        = "terraform-peer"
  port        = 0
}

resource "cloudflare_dns_zone_transfers_peer" "terraform_managed_resource_3" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  ixfr_enable = false
  name        = "fcusoidvut"
  port        = 0
}

