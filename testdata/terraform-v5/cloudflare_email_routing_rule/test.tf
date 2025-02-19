resource "cloudflare_email_routing_rule" "terraform_managed_resource" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  actions = [{
    type = "drop"
    value = ["destinationaddress@example.net"]
  }]
  matchers = [{
    field = "to"
    type = "literal"
    value = "test@example.com"
  }]
  enabled = true
  name = "Send to user@example.net rule."
  priority = 0
}
