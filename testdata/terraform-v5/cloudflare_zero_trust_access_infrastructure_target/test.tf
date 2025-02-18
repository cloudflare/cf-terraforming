resource "cloudflare_zero_trust_access_infrastructure_target" "example_zero_trust_access_infrastructure_target" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  hostname = "infra-access-target"
  ip = {
    ipv4 = {
      ip_addr = "187.26.29.249"
      virtual_network_id = "c77b744e-acc8-428f-9257-6878c046ed55"
    }
    ipv6 = {
      ip_addr = "64c0:64e8:f0b4:8dbf:7104:72b0:ec8f:f5e0"
      virtual_network_id = "c77b744e-acc8-428f-9257-6878c046ed55"
    }
  }
}
