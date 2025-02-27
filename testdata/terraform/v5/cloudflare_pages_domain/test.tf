resource "cloudflare_pages_domain" "terraform_managed_resource" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  project_name = "this-is-my-project-01"
  name = "example.com"
}
