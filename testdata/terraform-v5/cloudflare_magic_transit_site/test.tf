resource "cloudflare_magic_transit_site" "example_magic_transit_site" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  name = "site_1"
  connector_id = "ac60d3d0435248289d446cedd870bcf4"
  description = "description"
  ha_mode = true
  location = {
    lat = "37.6192"
    lon = "122.3816"
  }
  secondary_connector_id = "8d67040d3835dbcf46ce29da440dc482"
}
