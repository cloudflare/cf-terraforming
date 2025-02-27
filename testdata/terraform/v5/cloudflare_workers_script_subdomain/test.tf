resource "cloudflare_workers_script_subdomain" "terraform_managed_resource" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  script_name = "this-is_my_script-01"
  enabled = true
  previews_enabled = true
}
