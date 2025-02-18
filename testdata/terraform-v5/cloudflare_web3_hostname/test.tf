resource "cloudflare_web3_hostname" "example_web3_hostname" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  name = "gateway.example.com"
  target = "ethereum"
  description = "This is my IPFS gateway."
  dnslink = "/ipns/onboarding.ipfs.cloudflare.com"
}
