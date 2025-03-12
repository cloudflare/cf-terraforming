resource "cloudflare_email_routing_catch_all" "terraform_managed_resource" {
  enabled = false
  name    = "terraform rule catch all"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  actions = [{
    type  = "forward"
    value = ["destinationaddress@example.net"]
  }]
  matchers = [{
    type = "all"
  }]
}