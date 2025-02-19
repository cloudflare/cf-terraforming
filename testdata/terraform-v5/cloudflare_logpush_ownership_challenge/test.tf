resource "cloudflare_logpush_ownership_challenge" "terraform_managed_resource" {
  destination_conf = "s3://mybucket/logs?region=us-west-2"
  zone_id = "zone_id"
}
