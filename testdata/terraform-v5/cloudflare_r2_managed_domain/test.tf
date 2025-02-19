resource "cloudflare_r2_managed_domain" "terraform_managed_resource" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  bucket_name = "example-bucket"
  enabled = true
}
