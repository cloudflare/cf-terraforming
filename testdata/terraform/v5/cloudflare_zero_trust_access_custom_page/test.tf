resource "cloudflare_zero_trust_access_custom_page" "terraform_managed_resource" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  app_count   = 0
  custom_html = "<html><body><h1>Access Denied</h1></body></html>"
  name        = "plabknfrou"
  type        = "forbidden"
}

