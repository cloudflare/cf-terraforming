resource "cloudflare_custom_pages" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  state = "customized"
  type = "basic_challenge"
  url = "https://example.com/challenge.html"
}
