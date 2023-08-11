resource "cloudflare_ruleset" "terraform_managed_resource" {
  kind    = "zone"
  name    = "default"
  phase   = "http_request_late_transform"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  rules {
    action = "rewrite"
    action_parameters {
      headers {
        name      = "example-http-header-1"
        operation = "remove"
      }
      headers {
        name      = "example-http-header-2"
        operation = "remove"
      }
      headers {
        expression = "(ip.geoip.continent eq \"pluto\")"
        name       = "example-http-header-3"
        operation  = "set"
      }
      uri {
        path {
          value = "/aquarii_b"
        }
      }
    }
    description = "test transform"
    enabled     = true
    expression  = "(http.request.uri.path eq \"example.com\")"
  }
  rules {
    action = "rewrite"
    action_parameters {
      headers {
        name      = "example-http-static-header-1"
        operation = "set"
        value     = "my-http-header-1"
      }
    }
    description = "test transform set"
    enabled     = true
    expression  = "(http.request.uri.path eq \"example.com\")"
  }
  rules {
    action = "rewrite"
    action_parameters {
      uri {
        path {
          value = "/spaceship"
        }
      }
    }
    description = "test uri rewrite set"
    enabled     = false
    expression  = "(http.request.uri.path eq \"pumpkin.com\")"
  }
}
