resource "cloudflare_r2_managed_domain" "example_r2_managed_domain" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  bucket_name = "example-bucket"
  enabled = true
}
