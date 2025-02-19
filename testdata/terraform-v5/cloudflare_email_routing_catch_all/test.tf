resource "cloudflare_email_routing_catch_all" "terraform_managed_resource" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  actions = [{
    type = "drop"
    value = ["destinationaddress@example.net"]
  }]
  matchers = [{
    type = "all"
  }]
  enabled = true
  name = "Send to user@example.net rule."
}
