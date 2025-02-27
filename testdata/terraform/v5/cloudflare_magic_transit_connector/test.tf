resource "cloudflare_magic_transit_connector" "terraform_managed_resource" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  connector_id = "connector_id"
  activated = true
  interrupt_window_duration_hours = 0
  interrupt_window_hour_of_day = 0
  notes = "notes"
  timezone = "timezone"
}
