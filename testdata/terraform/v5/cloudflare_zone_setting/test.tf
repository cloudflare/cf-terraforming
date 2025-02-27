resource "cloudflare_zone_setting" "terraform_managed_resource" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  setting_id = "always_online"
  id = "0rtt"
  value = "on"
}
