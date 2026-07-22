resource "cloudflare_waiting_room_rules" "terraform_managed_resource" {
  waiting_room_id = "699d98642c564d2e855e9661899b7252"
  zone_id         = "0da42c8d2132a9ddaf714f9e7c920711"
  rules {
    action      = "bypass_waiting_room"
    description = "allow all traffic from 10.20.30.40"
    expression  = "ip.src in {10.20.30.40}"
  }
  rules {
    action      = "bypass_waiting_room"
    description = "allow all traffic from 50.60.70.80"
    expression  = "ip.src in {50.60.70.80}"
  }
}