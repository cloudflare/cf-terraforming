resource "cloudflare_stream_live_input" "example_stream_live_input" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  default_creator = "defaultCreator"
  delete_recording_after_days = 45
  meta = {
    name = "test stream 1"
  }
  recording = {
    allowed_origins = ["example.com"]
    hide_live_viewer_count = false
    mode = "off"
    require_signed_urls = false
    timeout_seconds = 0
  }
}
