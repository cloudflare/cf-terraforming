resource "cloudflare_account" "terraform_managed_resource_0" {
  name = "Foo Production"
  type = "enterprise"
  settings = {
    abuse_contact_email              = null
    access_approval_expiry           = null
    api_access_enabled               = null
    default_nameservers              = "cloudflare.standard"
    enforce_twofactor                = true
    use_account_custom_ns_by_default = false
  }
}

resource "cloudflare_account" "terraform_managed_resource_1" {
  name = "Foo Staging"
  type = "enterprise"
  settings = {
    abuse_contact_email              = null
    access_approval_expiry           = null
    api_access_enabled               = null
    default_nameservers              = "cloudflare.standard"
    enforce_twofactor                = false
    use_account_custom_ns_by_default = false
  }
}

resource "cloudflare_account" "terraform_managed_resource_2" {
  name = "Foo Acceptance Testing"
  type = "enterprise"
  settings = {
    abuse_contact_email              = null
    access_approval_expiry           = null
    api_access_enabled               = null
    default_nameservers              = "cloudflare.standard"
    enforce_twofactor                = false
    use_account_custom_ns_by_default = false
  }
}

