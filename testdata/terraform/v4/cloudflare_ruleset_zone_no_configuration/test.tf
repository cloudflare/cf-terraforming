resource "cloudflare_ruleset" "terraform_managed_resource" {
  description = "Some ruleset"
  kind        = "zone"
  name        = "default"
  phase       = "http_request_late_transform"
  zone_id     = "0da42c8d2132a9ddaf714f9e7c920711"
}
