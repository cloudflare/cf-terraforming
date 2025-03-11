resource "cloudflare_dns_zone_transfers_tsig" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  algo       = "hmac-sha512."
  name       = "fazhgkukxs."
  secret     = "caf79a7804b04337c9c66ccd7bef9190a1e1679b5dd03d8aa10f7ad45e1a9dab92b417896c15d4d007c7c14194538d2a5d0feffdecc5a7f0e1c570cfa700837c"
}

