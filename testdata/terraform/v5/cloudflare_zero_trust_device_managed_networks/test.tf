resource "cloudflare_zero_trust_device_managed_networks" "terraform_managed_resource_0" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  name       = "ibpcelbfgw"
  type       = "tls"
  config = {
    sha256       = "b5bb9d8014a0f9b1d61e21e796d78dccdf1352f23cd32812f4850b878ae4944c"
    tls_sockaddr = "foobar:1234"
  }
}

resource "cloudflare_zero_trust_device_managed_networks" "terraform_managed_resource_1" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  name       = "idivalyfsp"
  type       = "tls"
  config = {
    sha256       = "b5bb9d8014a0f9b1d61e21e796d78dccdf1352f23cd32812f4850b878ae4944c"
    tls_sockaddr = "foobar:1234"
  }
}

