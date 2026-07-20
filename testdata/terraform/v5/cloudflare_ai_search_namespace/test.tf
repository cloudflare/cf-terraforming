resource "cloudflare_ai_search_namespace" "terraform_managed_resource" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  description = "Test namespace"
  name        = "my-namespace"
}
