resource "cloudflare_waiting_room_rules" "example_waiting_room_rules" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  waiting_room_id = "699d98642c564d2e855e9661899b7252"
  rules = [{
    action = "bypass_waiting_room"
    expression = "ip.src in {10.20.30.40}"
    description = "allow all traffic from 10.20.30.40"
    enabled = true
  }]
}
