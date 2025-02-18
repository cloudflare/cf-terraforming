resource "cloudflare_logpush_ownership_challenge" "example_logpush_ownership_challenge" {
  destination_conf = "s3://mybucket/logs?region=us-west-2"
  zone_id = "zone_id"
}
