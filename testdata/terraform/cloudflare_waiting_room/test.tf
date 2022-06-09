resource "cloudflare_waiting_room" "terraform_managed_resource" {
  custom_page_html        = "{{#waitTimeKnown}} {{waitTime}} mins {{/waitTimeKnown}} {{^waitTimeKnown}} Queue all enabled {{/waitTimeKnown}}"
  description             = "Production - DO NOT MODIFY"
  disable_session_renewal = false
  host                    = "shop.example.com"
  json_response_enabled   = false
  name                    = "production_webinar"
  new_users_per_minute    = 1000
  path                    = "/shop/checkout"
  queue_all               = true
  session_duration        = 10
  suspended               = false
  total_active_users      = 1000
  zone_id                 = "0da42c8d2132a9ddaf714f9e7c920711"
}
