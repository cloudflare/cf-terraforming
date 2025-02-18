resource "cloudflare_dns_record" "example_dns_record" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  comment = "Domain verification record"
  content = "198.51.100.4"
  name = "example.com"
  proxied = true
  settings = {
    ipv4_only = true
    ipv6_only = true
  }
  tags = ["owner:dns-team"]
  ttl = 3600
  type = "A"
}
