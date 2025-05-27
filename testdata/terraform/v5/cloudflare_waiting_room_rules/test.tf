resource "cloudflare_waiting_room_rules" "terraform_managed_resource" {
  waiting_room_id = "8bbd1b13450f6c63ab6ab4e08a63762d"
  zone_id         = "0da42c8d2132a9ddaf714f9e7c920711"
  rules = [{
    action       = "bypass_waiting_room"
    description  = "cf-test"
    enabled      = true
    expression   = "(http.cookie eq \"foo\")"
    id           = "c5c159572b7a44a78bffd87ac2d6457d"
    last_updated = "2025-05-27T18:26:07.047916Z"
    version      = "1"
    }, {
    action       = "bypass_waiting_room"
    description  = "rule-2"
    enabled      = true
    expression   = "(ip.src.is_in_european_union)"
    id           = "30bb12976d124f2aacb2335cda5b2817"
    last_updated = "2025-05-27T18:53:44.853226Z"
    version      = "1"
  }]
}

