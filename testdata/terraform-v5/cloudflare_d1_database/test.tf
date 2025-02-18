resource "cloudflare_d1_database" "example_d1_database" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  name = "my-database"
  primary_location_hint = "wnam"
}
