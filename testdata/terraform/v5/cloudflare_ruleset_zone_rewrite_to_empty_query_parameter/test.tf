resource "cloudflare_ruleset" "terraform_managed_resource" {
  kind    = "zone"
  name    = "default"
  phase   = "http_request_transform"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  rules = [{
    action = "rewrite"
    action_parameters = {
      uri = {
        query = {
          value = ""
        }
      }
    }
    description  = "rewrite with no query string"
    enabled      = true
    expression   = "true"
    id           = null
    last_updated = "2023-02-16T00:26:08.978517Z"
    ref          = "1fb6a3117e864d46bcda192d14a1e1dc"
    version      = "5"
  }]
}

