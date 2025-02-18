resource "cloudflare_dns_zone_transfers_incoming" "example_dns_zone_transfers_incoming" {
  zone_id = "269d8f4853475ca241c4e730be286b20"
  auto_refresh_seconds = 86400
  name = "www.example.com."
  peers = ["23ff594956f20c2a721606e94745a8aa", "00920f38ce07c2e2f4df50b1f61d4194"]
}
