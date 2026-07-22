resource "cloudflare_stream_live_input" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  meta = jsonencode({
    name = "test stream 1"
  })
}

