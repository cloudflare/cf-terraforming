resource "cloudflare_waiting_room" "terraform_managed_resource" {
  zone_id                 = "0da42c8d2132a9ddaf714f9e7c920711"
  name                    = "production_webinar"
  description             = "Production - DO NOT MODIFY"
  suspended               = false
  host                    = "shop.example.com"
  path                    = "/shop/checkout"
  queue_all               = true
  new_users_per_minute    = 1000
  total_active_users      = 1000
  session_duration        = 10
  disable_session_renewal = false
  json_response_enabled   = false
  custom_page_html        = "{{#waitTimeKnown}} {{waitTime}} mins {{/waitTimeKnown}} {{^waitTimeKnown}} Queue all enabled {{/waitTimeKnown}}"
}
