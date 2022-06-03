resource "cloudflare_page_rule" "terraform_managed_resource" {
  priority = 1
  status   = "active"
  target   = "*example.com/images/*"
  zone_id  = "0da42c8d2132a9ddaf714f9e7c920711"
  actions {
    cache_key_fields {
      cookie {
        check_presence = ["x-some-header"]
        include        = ["x-some-other-header", "x-some-value"]
      }
      header {
        check_presence = ["x-forwarded-for"]
        include        = ["authorization"]
      }
      host {
        resolved = true
      }
      query_string {
        exclude = ["*"]
      }
      user {
        device_type = true
        geo         = true
        lang        = false
      }
    }
    cache_ttl_by_status {
      codes = "200"
      ttl   = 60
    }
    cache_ttl_by_status {
      codes = "204"
      ttl   = -1
    }
    cache_ttl_by_status {
      codes = "300-399"
      ttl   = 0
    }
    browser_cache_ttl    = 1800
    disable_apps         = true
    host_header_override = "not-example.com"
  }
}
