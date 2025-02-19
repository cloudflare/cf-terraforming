resource "cloudflare_api_shield_schema_validation_settings" "terraform_managed_resource" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  validation_default_mitigation_action = "none"
  validation_override_mitigation_action = "none"
}
