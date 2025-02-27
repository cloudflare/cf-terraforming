resource "cloudflare_turnstile_widget" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  domains    = ["example.com", "example.net"]
  mode       = "non-interactive"
  name       = "example sitekey"
  region     = "world"
}
