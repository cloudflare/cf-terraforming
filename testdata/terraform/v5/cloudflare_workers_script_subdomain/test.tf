resource "cloudflare_workers_script_subdomain" "terraform_managed_resource" {
  account_id       = "f037e56e89293a057740de681ac9abbe"
  enabled          = true
  previews_enabled = true
  script_name      = "accounts"
}

