resource "cloudflare_zone_setting" "terraform_managed_resource" {
  setting_id = "always_online"
  zone_id    = "0da42c8d2132a9ddaf714f9e7c920711"
  value      = "off"
}
