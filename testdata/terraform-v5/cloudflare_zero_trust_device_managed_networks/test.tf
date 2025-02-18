resource "cloudflare_zero_trust_device_managed_networks" "example_zero_trust_device_managed_networks" {
  account_id = "699d98642c564d2e855e9661899b7252"
  config = {
    tls_sockaddr = "foobar:1234"
    sha256 = "b5bb9d8014a0f9b1d61e21e796d78dccdf1352f23cd32812f4850b878ae4944c"
  }
  name = "managed-network-1"
  type = "tls"
}
