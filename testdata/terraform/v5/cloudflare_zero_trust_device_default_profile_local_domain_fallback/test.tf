resource "cloudflare_zero_trust_device_default_profile_local_domain_fallback" "terraform_managed_resource_0" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  domains = [{
    description = "example domain"
    dns_server  = ["1.0.0.1"]
    suffix      = "example.com"
  }]
}

resource "cloudflare_zero_trust_device_default_profile_local_domain_fallback" "terraform_managed_resource_1" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  domains = [{
    description = "foo domain"
    dns_server  = ["1.0.0.1"]
    suffix      = "foo.com"
  }]
}

