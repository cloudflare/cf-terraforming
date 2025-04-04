resource "cloudflare_zero_trust_dlp_dataset" "terraform_managed_resource" {
  account_id       = "f037e56e89293a057740de681ac9abbe"
  encoding_version = 1
  name             = "tf-test"
  secret           = true
}

