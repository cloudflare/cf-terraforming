resource "cloudflare_zero_trust_access_ai_controls_mcp_portal" "terraform_managed_resource" {
  account_id         = "f037e56e89293a057740de681ac9abbe"
  description        = "A test MCP portal for cf-terraforming"
  hostname           = "mcp-portal.terraform.cfapi.net"
  name               = "test-mcp-portal"
  secure_web_gateway = false
  servers = [{
    auth_type        = "unauthenticated"
    created_at       = "2026-03-11 15:00:45"
    created_by       = "tamas@cloudflare.com"
    default_disabled = true
    description      = "A test MCP server for cf-terraforming"
    error            = "Server error: HTTP 530"
    hostname         = "https://mcp.example.com"
    id               = "a1b2c3d4e5f67890abcdef1234567890"
    last_synced      = "2026-03-11 15:00:45"
    modified_at      = "2026-03-11 15:00:45"
    modified_by      = "tamas@cloudflare.com"
    name             = "test-mcp-server"
    on_behalf        = true
    prompts          = []
    status           = "error"
    tools            = []
    updated_prompts = [{
      description = "System context prompt"
      enabled     = true
      name        = "system"
      }, {
      description = "Initial greeting"
      enabled     = true
      name        = "greeting"
    }]
    updated_tools = [{
      description = "Search for information"
      enabled     = true
      name        = "search"
      }, {
      description = "Summarize content"
      enabled     = true
      name        = "summarize"
      }, {
      description = "Translate text"
      enabled     = false
      name        = "translate"
    }]
  }]
}
