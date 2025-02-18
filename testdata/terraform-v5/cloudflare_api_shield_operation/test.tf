resource "cloudflare_api_shield_operation" "example_api_shield_operation" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  endpoint = "/api/v1/users/{var1}"
  host = "www.example.com"
  method = "GET"
}
