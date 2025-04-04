resource "cloudflare_zero_trust_dns_location" "terraform_managed_resource" {
  account_id             = "f037e56e89293a057740de681ac9abbe"
  client_default         = true
  dns_destination_ips_id = "0e4a32c6-6fb8-4858-9296-98f51631e8e6"
  ecs_support            = false
  name                   = "bbltsfpzao"
  networks               = []
}

