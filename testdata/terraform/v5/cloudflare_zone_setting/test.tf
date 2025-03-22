resource "cloudflare_zone_setting" "terraform_managed_resource_0" {
  setting_id = "always_online"
  zone_id    = "0da42c8d2132a9ddaf714f9e7c920711"
  value      = "off"
}

resource "cloudflare_zone_setting" "terraform_managed_resource_1" {
  setting_id = "cache_level"
  zone_id    = "0da42c8d2132a9ddaf714f9e7c920711"
  value      = "aggressive"
}

