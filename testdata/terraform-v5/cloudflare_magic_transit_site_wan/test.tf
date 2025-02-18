resource "cloudflare_magic_transit_site_wan" "example_magic_transit_site_wan" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  site_id = "023e105f4ecef8ad9ca31a8372d0c353"
  physport = 1
  vlan_tag = 0
  name = "name"
  priority = 0
  static_addressing = {
    address = "192.0.2.0/24"
    gateway_address = "192.0.2.1"
    secondary_address = "192.0.2.0/24"
  }
}
