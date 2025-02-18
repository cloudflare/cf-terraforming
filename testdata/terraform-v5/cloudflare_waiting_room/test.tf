resource "cloudflare_waiting_room" "example_waiting_room" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  host = "shop.example.com"
  name = "production_webinar"
  new_users_per_minute = 200
  total_active_users = 200
  additional_routes = [{
    host = "shop2.example.com"
    path = "/shop2/checkout"
  }]
  cookie_attributes = {
    samesite = "auto"
    secure = "auto"
  }
  cookie_suffix = "abcd"
  custom_page_html = "{{#waitTimeKnown}} {{waitTime}} mins {{/waitTimeKnown}} {{^waitTimeKnown}} Queue all enabled {{/waitTimeKnown}}"
  default_template_language = "en-US"
  description = "Production - DO NOT MODIFY"
  disable_session_renewal = false
  enabled_origin_commands = ["revoke"]
  json_response_enabled = false
  path = "/shop/checkout"
  queue_all = true
  queueing_method = "fifo"
  queueing_status_code = 200
  session_duration = 1
  suspended = true
  turnstile_action = "log"
  turnstile_mode = "off"
}
