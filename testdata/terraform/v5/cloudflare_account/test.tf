resource "cloudflare_account" "terraform_managed_resource_0" {
  name = "Foo Production"
  type = "enterprise"
  settings = {
    default_nameservers = "cloudflare.standard"
  }
}

resource "cloudflare_account" "terraform_managed_resource_1" {
  name = "Foo Staging"
  type = "enterprise"
  settings = {
    default_nameservers = "cloudflare.standard"
  }
}

resource "cloudflare_account" "terraform_managed_resource_2" {
  name = "Foo Acceptance Testing"
  type = "enterprise"
  settings = {
    default_nameservers = "cloudflare.standard"
  }
}