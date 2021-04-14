resource "cloudflare_access_identity_provider" "terraform_managed_resource" {
  name = "PIN login"
  type = "onetimepin"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
}
