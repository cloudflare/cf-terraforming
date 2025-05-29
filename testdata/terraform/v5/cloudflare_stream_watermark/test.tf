resource "cloudflare_stream_watermark" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  file       = "REPLACE with filebase64(\"path-to-file\")"
  name       = "Marketing Videos"
  opacity    = 0.75
  padding    = 0.1
  position   = "center"
  scale      = 0.1
}

