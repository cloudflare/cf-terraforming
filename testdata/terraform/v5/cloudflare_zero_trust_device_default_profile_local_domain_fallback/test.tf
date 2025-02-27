resource "cloudflare_zero_trust_device_default_profile_local_domain_fallback" "terraform_managed_resource" {
  account_id = "699d98642c564d2e855e9661899b7252"
  domains = [{
    suffix = "example.com"
    description = "Domain bypass for local development"
    dns_server = ["1.1.1.1"]
  }]
}
