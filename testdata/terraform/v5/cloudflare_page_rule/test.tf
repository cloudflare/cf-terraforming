resource "cloudflare_page_rule" "terraform_managed_resource_0" {
  priority = 2
  status   = "active"
  target   = "*terraform.cfapi.net/_assets/*"
  zone_id  = "0da42c8d2132a9ddaf714f9e7c920711"
  actions = {
    cache_level            = "cache_everything"
    explicit_cache_control = "on"
    host_header_override   = "prod-notion-assets.terraform.cfapi.net"
    resolve_override       = "prod-notion-assets-s3.terraform.cfapi.net"
    ssl                    = "full"
  }
}

resource "cloudflare_page_rule" "terraform_managed_resource_1" {
  priority = 1
  status   = "active"
  target   = "www.example./com/*"
  zone_id  = "0da42c8d2132a9ddaf714f9e7c920711"
  actions = {
    email_obfuscation = "off"
  }
}

