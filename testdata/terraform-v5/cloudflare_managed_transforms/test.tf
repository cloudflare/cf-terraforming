resource "cloudflare_managed_transforms" "example_managed_transforms" {
  zone_id = "9f1839b6152d298aca64c4e906b6d074"
  managed_request_headers = [{
    id = "add_bot_protection_headers"
    enabled = true
  }]
  managed_response_headers = [{
    id = "add_security_headers"
    enabled = true
  }]
}
