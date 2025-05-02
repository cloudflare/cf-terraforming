resource "cloudflare_waiting_room_event" "terraform_managed_resource" {
  custom_page_html       = "{{#waitTimeKnown}} {{waitTime}} mins {{/waitTimeKnown}} {{^waitTimeKnown}} Event is prequeueing / Queue all enabled {{/waitTimeKnown}}"
  description            = "Production event - DO NOT MODIFY"
  event_end_time         = "2021-09-28T17:00:00Z"
  event_start_time       = "2021-09-28T15:30:00Z"
  name                   = "production_webinar_event"
  prequeue_start_time    = "2021-09-28T15:00:00Z"
  queueing_method        = "random"
  shuffle_at_event_start = false
  suspended              = false
  waiting_room_id        = "e7f9e4c190ea8d6c66cab32ac110f39a"
  zone_id                = "0da42c8d2132a9ddaf714f9e7c920711"
}

