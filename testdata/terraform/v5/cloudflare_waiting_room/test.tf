resource "cloudflare_waiting_room" "terraform_managed_resource" {
  custom_page_html          = "foobar"
  default_template_language = "en-US"
  description               = "my desc"
  disable_session_renewal   = true
  host                      = "www.terraform.cfapi.net"
  json_response_enabled     = true
  name                      = "waiting_room_ucmxvksksg"
  new_users_per_minute      = 400
  path                      = "/foobar"
  queue_all                 = false
  queueing_method           = "fifo"
  queueing_status_code      = 200
  session_duration          = 10
  suspended                 = true
  total_active_users        = 405
  turnstile_action          = "log"
  turnstile_mode            = "invisible"
  zone_id                   = "0da42c8d2132a9ddaf714f9e7c920711"
  cookie_attributes = {
    samesite = "auto"
    secure   = "auto"
  }
}

