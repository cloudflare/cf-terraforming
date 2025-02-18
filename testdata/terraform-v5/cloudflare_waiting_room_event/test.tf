resource "cloudflare_waiting_room_event" "example_waiting_room_event" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  waiting_room_id = "699d98642c564d2e855e9661899b7252"
  event_end_time = "2021-09-28T17:00:00.000Z"
  event_start_time = "2021-09-28T15:30:00.000Z"
  name = "production_webinar_event"
  custom_page_html = "{{#waitTimeKnown}} {{waitTime}} mins {{/waitTimeKnown}} {{^waitTimeKnown}} Event is prequeueing / Queue all enabled {{/waitTimeKnown}}"
  description = "Production event - DO NOT MODIFY"
  disable_session_renewal = true
  new_users_per_minute = 200
  prequeue_start_time = "2021-09-28T15:00:00.000Z"
  queueing_method = "random"
  session_duration = 1
  shuffle_at_event_start = true
  suspended = true
  total_active_users = 200
}
