resource "cloudflare_zero_trust_access_infrastructure_target" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  hostname   = "infra-access-target"
  ip = {
    ipv4 = {
      ip_addr            = "187.26.29.249"
      virtual_network_id = "59c65fed-41cd-4d00-a861-a1bd3b90a32f"
    }
    ipv6 = null
  }
}

