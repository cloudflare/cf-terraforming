resource "cloudflare_custom_ssl" "terraform_managed_resource" {
  custom_ssl_options {
    bundle_method    = "ubiquitous"
    certificate      = "-----INSERT CERTIFICATE-----"
    geo_restrictions = "us"
    private_key      = "-----INSERT PRIVATE KEY-----"
    type             = "legacy_custom"
  }
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
}
