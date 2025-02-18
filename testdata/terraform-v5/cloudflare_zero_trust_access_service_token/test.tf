resource "cloudflare_zero_trust_access_service_token" "example_zero_trust_access_service_token" {
  name = "CI/CD token"
  zone_id = "zone_id"
  duration = "60m"
}
