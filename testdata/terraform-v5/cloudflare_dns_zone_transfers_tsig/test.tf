resource "cloudflare_dns_zone_transfers_tsig" "example_dns_zone_transfers_tsig" {
  account_id = "01a7362d577a6c3019a474fd6f485823"
  algo = "hmac-sha512."
  name = "tsig.customer.cf."
  secret = "caf79a7804b04337c9c66ccd7bef9190a1e1679b5dd03d8aa10f7ad45e1a9dab92b417896c15d4d007c7c14194538d2a5d0feffdecc5a7f0e1c570cfa700837c"
}
