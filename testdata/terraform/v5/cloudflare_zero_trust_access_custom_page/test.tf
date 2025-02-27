resource "cloudflare_zero_trust_access_custom_page" "terraform_managed_resource" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  custom_html = "<html><body><h1>Access Denied</h1></body></html>"
  name = "name"
  type = "identity_denied"
  app_count = 0
}
