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
        value      = "space-header"
      }
      uri {
        path {
          value = "/aquarii_b"
        }
      }
    }
    description  = "test transform"
    enabled      = true
    expression   = "(http.request.uri.path eq \"example.com\")"
    id           = "e5b61605d6cf4ce08f729c17d42d76ef"
    last_updated = "2022-02-07T16:58:54.317608Z"
    ref          = "e5b61605d6cf4ce08f729c17d42d76ef"
    version      = "1"
  }
  rules {
    action = "rewrite"
    action_parameters {
      headers {
        expression = "(ip.geoip.continent eq \"T1\")"
        name       = "example-http-static-header-1"
        operation  = "set"
        value      = "my-http-header-1"
      }
    }
    description  = "test transform set"
    enabled      = true
    expression   = "(http.request.uri.path eq \"example.com\")"
    id           = "8ec764cf386940c89dd83dbab7bb4c16"
    last_updated = "2022-02-07T16:58:54.317608Z"
    ref          = "8ec764cf386940c89dd83dbab7bb4c16"
    version      = "1"
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
    description  = "test uri rewrite set"
    enabled      = false
    expression   = "(http.request.uri.path eq \"pumpkin.com\")"
    id           = "d0f1b4fdb4234adf9c6de9b614424836"
    last_updated = "2022-05-07T16:58:54.317608Z"
    ref          = "d0f1b4fdb4234adf9c6de9b614424836"
    version      = "1"
  }
}
