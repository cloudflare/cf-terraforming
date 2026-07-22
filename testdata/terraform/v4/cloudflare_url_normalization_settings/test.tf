resource "cloudflare_url_normalization_settings" "terraform_managed_resource" {
  scope   = "incoming"
  type    = "cloudflare"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
}
