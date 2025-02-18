resource "cloudflare_zone_dnssec" "example_zone_dnssec" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  dnssec_multi_signer = false
  dnssec_presigned = true
  status = "active"
}
