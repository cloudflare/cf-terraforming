resource "cloudflare_email_routing_rule" "terraform_managed_resource_0" {
  enabled  = true
  name     = "Rule created at 2025-02-21T21: 29: 30.992Z"
  priority = 0
  zone_id  = "0da42c8d2132a9ddaf714f9e7c920711"
  actions = [{
    type = "drop"
  }]
  matchers = [{
    field = "to"
    type  = "literal"
    value = "abcd@example.com"
  }]
}

resource "cloudflare_email_routing_rule" "terraform_managed_resource_1" {
  enabled  = false
  name     = "terraform rule catch all"
  priority = 2147483647
  zone_id  = "0da42c8d2132a9ddaf714f9e7c920711"
  actions = [{
    type  = "forward"
    value = ["destinationaddress@example.net"]
  }]
  matchers = [{
    type = "all"
  }]
}

