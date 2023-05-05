resource "cloudflare_tiered_cache" "terraform_managed_resource" {
  cache_type              = "smart"
  zone_id                 = "0da42c8d2132a9ddaf714f9e7c920711"
}
