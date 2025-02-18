resource "cloudflare_hyperdrive_config" "example_hyperdrive_config" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  name = "example-hyperdrive"
  origin = {
    database = "postgres"
    host = "database.example.com"
    password = "password"
    port = 5432
    scheme = "postgres"
    user = "postgres"
  }
  caching = {
    disabled = true
  }
}
