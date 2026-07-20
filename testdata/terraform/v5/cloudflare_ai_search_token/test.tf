resource "cloudflare_ai_search_token" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  cf_api_id  = "api-id-123"
  cf_api_key = "api-key-secret"
  legacy     = true
  name       = "my-token"
}
