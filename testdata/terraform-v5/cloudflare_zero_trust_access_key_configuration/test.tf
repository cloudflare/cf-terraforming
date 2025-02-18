resource "cloudflare_zero_trust_access_key_configuration" "example_zero_trust_access_key_configuration" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  key_rotation_interval_days = 30
}
