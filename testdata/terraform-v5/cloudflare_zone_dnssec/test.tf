resource "cloudflare_zone_dnssec" "terraform_managed_resource" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  dnssec_multi_signer = false
  dnssec_presigned = true
  status = "active"
}
