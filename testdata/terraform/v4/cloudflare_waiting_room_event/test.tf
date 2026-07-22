resource "cloudflare_waiting_room_event" "terraform_managed_resource" {
  custom_page_html        = "{{#waitTimeKnown}} {{waitTime}} mins {{/waitTimeKnown}} {{^waitTimeKnown}} Event is prequeueing / Queue all enabled {{/waitTimeKnown}}"
  description             = "Production event - DO NOT MODIFY"
  disable_session_renewal = false
  event_end_time          = "2021-09-28T17:00:00Z"
  event_start_time        = "2021-09-28T15:30:00Z"
  name                    = "production_webinar_event"
  new_users_per_minute    = 1000
  queueing_method         = "fifo"
  session_duration        = 10
  shuffle_at_event_start  = false
  suspended               = false
  total_active_users      = 1000
  waiting_room_id         = "699d98642c564d2e855e9661899b7252"
  zone_id                 = "0da42c8d2132a9ddaf714f9e7c920711"
}
