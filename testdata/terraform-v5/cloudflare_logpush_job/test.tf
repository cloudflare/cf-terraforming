resource "cloudflare_logpush_job" "example_logpush_job" {
  destination_conf = "s3://mybucket/logs?region=us-west-2"
  zone_id = "zone_id"
  dataset = "http_requests"
  enabled = false
  frequency = "high"
  kind = "edge"
  logpull_options = "fields=RayID,ClientIP,EdgeStartTimestamp&timestamps=rfc3339"
  max_upload_bytes = 5000000
  max_upload_interval_seconds = 30
  max_upload_records = 1000
  name = "example.com"
  output_options = {
    batch_prefix = "batch_prefix"
    batch_suffix = "batch_suffix"
    cve_2021_4428 = true
    field_delimiter = "field_delimiter"
    field_names = ["ClientIP", "EdgeStartTimestamp", "RayID"]
    output_type = "ndjson"
    record_delimiter = "record_delimiter"
    record_prefix = "record_prefix"
    record_suffix = "record_suffix"
    record_template = "record_template"
    sample_rate = 0
    timestamp_format = "unixnano"
  }
  ownership_challenge = "00000000000000000000"
}
