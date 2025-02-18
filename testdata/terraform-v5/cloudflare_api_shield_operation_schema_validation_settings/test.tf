resource "cloudflare_api_shield_operation_schema_validation_settings" "example_api_shield_operation_schema_validation_settings" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  operation_id = "f174e90a-fafe-4643-bbbc-4a0ed4fc8415"
  mitigation_action = "log"
}
