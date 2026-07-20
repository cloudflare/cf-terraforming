resource "cloudflare_ai_gateway_dynamic_routing" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  name       = "my-route"
  elements = [{
    id      = "elem-1"
    outputs = {}
    properties = {
      ai_gateway_dynamic_routing_provider = "openai"
      model                               = "gpt-4"
    }
    type = "model"
  }]
}
