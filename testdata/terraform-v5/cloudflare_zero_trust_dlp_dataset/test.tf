resource "cloudflare_zero_trust_dlp_dataset" "terraform_managed_resource" {
  account_id = "account_id"
  name = "name"
  description = "description"
  encoding_version = 0
  secret = true
}
