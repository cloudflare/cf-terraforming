resource "cloudflare_r2_bucket" "example_r2_bucket" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  name = "example-bucket"
  location = "apac"
  storage_class = "Standard"
}
