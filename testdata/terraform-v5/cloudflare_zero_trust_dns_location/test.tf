resource "cloudflare_zero_trust_dns_location" "example_zero_trust_dns_location" {
  account_id = "699d98642c564d2e855e9661899b7252"
  name = "Austin Office Location"
  client_default = false
  dns_destination_ips_id = "0e4a32c6-6fb8-4858-9296-98f51631e8e6"
  ecs_support = false
  endpoints = {
    doh = {
      enabled = true
      networks = [{
        network = "2001:85a3::/64"
      }]
      require_token = true
    }
    dot = {
      enabled = true
      networks = [{
        network = "2001:85a3::/64"
      }]
    }
    ipv4 = {
      enabled = true
    }
    ipv6 = {
      enabled = true
      networks = [{
        network = "2001:85a3::/64"
      }]
    }
  }
  networks = [{
    network = "192.0.2.1/32"
  }]
}
