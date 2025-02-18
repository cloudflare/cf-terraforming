resource "cloudflare_api_shield_schema" "example_api_shield_schema" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  file = "file.txt"
  kind = "openapi_v3"
  name = "petstore schema"
  validation_enabled = "true"
}
