resource "cloudflare_api_shield_schema_validation_settings" "terraform_managed_resource" {
  validation_default_mitigation_action = "log"
  zone_id                              = "0da42c8d2132a9ddaf714f9e7c920711"
}

