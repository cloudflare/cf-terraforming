resource "cloudflare_d1_database" "terraform_managed_resource" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  name = "my-database"
  primary_location_hint = "wnam"
}
