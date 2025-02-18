resource "cloudflare_regional_hostname" "example_regional_hostname" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  hostname = "foo.example.com"
  region_key = "ca"
}
