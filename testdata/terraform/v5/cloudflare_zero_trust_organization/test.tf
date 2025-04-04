resource "cloudflare_zero_trust_organization" "terraform_managed_resource" {
  account_id                         = "f037e56e89293a057740de681ac9abbe"
  allow_authenticate_via_warp        = false
  auth_domain                        = "lklfsevdnw-terraform-cfapi.cloudflareaccess.com"
  is_ui_read_only                    = false
  name                               = "terraform-cfapi.cloudflareaccess.com"
  session_duration                   = "12h"
  user_seat_expiration_inactive_time = "1460h"
  warp_auth_session_duration         = "36h"
  login_design = {
    background_color = "#FFFFFF"
    footer_text      = "My footer text"
    header_text      = "My header text updated"
    logo_path        = "https://example.com/logo.png"
    text_color       = "#000000"
  }
}

