resource "cloudflare_logpush_job" "terraform_managed_resource" {
  dataset = "http_requests"
  destination_conf = "s3://mybucket/logs?region=us-west-2"
  enabled = false
  logpull_options = "fields=RayID,ClientIP,EdgeStartTimestamp&timestamps=rfc3339"
  name = "example.com"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
}
