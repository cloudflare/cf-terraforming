resource "cloudflare_zone" "terraform_managed_resource_0" {
  name                = "api.cloudflarecn.net"
  type                = "full"
  vanity_name_servers = []
  account = {
    id   = "d781e89412e1965584dc2a65a43b7fce"
    name = "Foo Team Production"
  }
}

resource "cloudflare_zone" "terraform_managed_resource_1" {
  name                = "api.cloudflare.com"
  type                = "full"
  vanity_name_servers = []
  account = {
    id   = "d781e89412e1965584dc2a65a43b7fce"
    name = "API Team Production"
  }
}

resource "cloudflare_zone" "terraform_managed_resource_2" {
  name                = "api.staging.cloudflarecn.net"
  type                = "full"
  vanity_name_servers = []
  account = {
    id   = "d98242917a12ae74b03c3e3058dbacbd"
    name = "Foo Team Staging"
  }
}

