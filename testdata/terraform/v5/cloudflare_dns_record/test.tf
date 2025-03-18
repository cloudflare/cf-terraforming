resource "cloudflare_dns_record" "terraform_managed_resource_0" {
  content  = "1.1.1.1"
  name     = "foo.example.com"
  proxied  = false
  tags     = []
  ttl      = 1
  type     = "A"
  zone_id  = "0da42c8d2132a9ddaf714f9e7c920711"
  settings = {}
}

resource "cloudflare_dns_record" "terraform_managed_resource_1" {
  content = "example.com"
  name    = "atmdfzvyns.origin.example.com"
  proxied = false
  tags    = []
  ttl     = 3600
  type    = "CNAME"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  settings = {
    flatten_cname = false
  }
}

resource "cloudflare_dns_record" "terraform_managed_resource_2" {
  content  = "mx.record.example.com"
  name     = "hwflxxxmoc.example.com"
  priority = 71
  proxied  = false
  tags     = []
  ttl      = 1
  type     = "MX"
  zone_id  = "0da42c8d2132a9ddaf714f9e7c920711"
  settings = {}
}

