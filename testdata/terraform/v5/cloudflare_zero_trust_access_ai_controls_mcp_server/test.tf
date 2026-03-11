resource "cloudflare_zero_trust_access_ai_controls_mcp_server" "terraform_managed_resource" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  auth_type   = "unauthenticated"
  description = "A test MCP server for cf-terraforming"
  hostname    = "https://mcp.example.com"
  name        = "test-mcp-server"
}
