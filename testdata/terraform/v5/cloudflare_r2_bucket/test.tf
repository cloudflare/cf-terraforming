resource "cloudflare_r2_bucket" "terraform_managed_resource" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  name = "example-bucket"
  location = "apac"
  storage_class = "Standard"
}
