resource "cloudflare_custom_ssl" "terraform_managed_resource" {
  bundle_method = "ubiquitous"
  policy        = "(country: US) or (region: EU)"
  zone_id       = "023e105f4ecef8ad9ca31a8372d0c353"
  geo_restrictions = {
    label = "us"
  }
}
