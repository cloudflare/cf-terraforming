resource "cloudflare_zero_trust_gateway_settings" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  settings = {
    activity_log = {
      enabled = true
    }
    antivirus = {
      enabled_download_phase = true
      enabled_upload_phase   = false
      fail_closed            = true
      notification_settings = {
        enabled     = true
        msg         = "msg"
        support_url = "https://hello.com/"
      }
    }
    block_page = {
      background_color = "#000000"
      enabled          = true
      footer_text      = "hello"
      header_text      = "hello"
      include_context  = false
      logo_path        = "https://example.com"
      mailto_address   = "test@cloudflare.com"
      mailto_subject   = "hello"
      name             = "iddghecuxq"
      suppress_footer  = false
      target_uri       = ""
    }
    body_scanning = {
      inspection_mode = "deep"
    }
    browser_isolation = {
      non_identity_enabled          = false
      url_browser_isolation_enabled = true
    }
    custom_certificate = {
      created_at  = "0001-01-01T00:00:00Z"
      enabled     = false
      id          = "00000000-0000-0000-0000-000000000000"
      qs_pack_id  = "00000000-0000-0000-0000-000000000000"
      uploaded_on = "0001-01-01T00:00:00Z"
    }
    extended_email_matching = {
      enabled = true
    }
    fips = {
      tls = true
    }
    protocol_detection = {
      enabled = true
    }
    tls_decrypt = {
      enabled = true
    }
  }
}

