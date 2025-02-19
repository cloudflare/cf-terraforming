resource "cloudflare_zero_trust_access_service_token" "terraform_managed_resource" {
  name = "CI/CD token"
  zone_id = "zone_id"
  duration = "60m"
}
