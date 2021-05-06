resource "cloudflare_page_rule" "terraform_managed_resource_9a7806061c88ada191ed06f989cc3dac" {
  priority = 1
  status = "active"
  target = "*example.com/images/*"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  actions {
    minify {
      css = "on"
      html = "on"
      js = "off"
    }
    always_online = "on"
  }
}