resource "cloudflare_magic_transit_site_lan" "example_magic_transit_site_lan" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  site_id = "023e105f4ecef8ad9ca31a8372d0c353"
  physport = 1
  vlan_tag = 0
  ha_link = true
  name = "name"
  nat = {
    static_prefix = "192.0.2.0/24"
  }
  routed_subnets = [{
    next_hop = "192.0.2.1"
    prefix = "192.0.2.0/24"
    nat = {
      static_prefix = "192.0.2.0/24"
    }
  }]
  static_addressing = {
    address = "192.0.2.0/24"
    dhcp_relay = {
      server_addresses = ["192.0.2.1"]
    }
    dhcp_server = {
      dhcp_pool_end = "192.0.2.1"
      dhcp_pool_start = "192.0.2.1"
      dns_server = "192.0.2.1"
      reservations = {
        "00:11:22:33:44:55" = "192.0.2.100"
        "AA:BB:CC:DD:EE:FF" = "192.168.1.101"
      }
    }
    secondary_address = "192.0.2.0/24"
    virtual_address = "192.0.2.0/24"
  }
}
