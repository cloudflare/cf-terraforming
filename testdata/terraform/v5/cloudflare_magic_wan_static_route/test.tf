resource "cloudflare_magic_wan_static_route" "terraform_managed_resource_0" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  description = "VLAN-125-IPSec2"
  nexthop     = "172.16.255.3"
  prefix      = "192.168.125.0/24"
  priority    = 100
}

resource "cloudflare_magic_wan_static_route" "terraform_managed_resource_1" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  description = "VLAN-125-IPSec1"
  nexthop     = "172.16.255.1"
  prefix      = "192.168.125.0/24"
  priority    = 100
}

