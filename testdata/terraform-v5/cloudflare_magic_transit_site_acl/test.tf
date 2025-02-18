resource "cloudflare_magic_transit_site_acl" "example_magic_transit_site_acl" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  site_id = "023e105f4ecef8ad9ca31a8372d0c353"
  lan_1 = {
    lan_id = "lan_id"
    lan_name = "lan_name"
    port_ranges = ["8080-9000"]
    ports = [1]
    subnets = ["192.0.2.1"]
  }
  lan_2 = {
    lan_id = "lan_id"
    lan_name = "lan_name"
    port_ranges = ["8080-9000"]
    ports = [1]
    subnets = ["192.0.2.1"]
  }
  name = "PIN Pad - Cash Register"
  description = "Allows local traffic between PIN pads and cash register."
  forward_locally = true
  protocols = ["tcp"]
  unidirectional = true
}
