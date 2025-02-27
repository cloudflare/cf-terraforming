resource "cloudflare_url_normalization_settings" "terraform_managed_resource" {
  zone_id = "9f1839b6152d298aca64c4e906b6d074"
  scope = "incoming"
  type = "cloudflare"
}
