resource "cloudflare_zone" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  zone       = "example.com"
}
