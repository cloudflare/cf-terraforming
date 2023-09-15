resource "cloudflare_access_identity_provider" "terraform_managed_resource" {
  name    = "GitHub OAuth"
  type    = "github"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  config {
    client_id = "example-id"
  }
}
