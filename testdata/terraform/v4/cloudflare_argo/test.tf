resource "cloudflare_argo" "terraform_managed_resource" {
  smart_routing  = "on"
  tiered_caching = "off"
  zone_id        = "0da42c8d2132a9ddaf714f9e7c920711"
}
