resource "cloudflare_ruleset" "terraform_managed_resource" {
  kind    = "zone"
  name    = "Zone sanitize ruleset"
  phase   = "http_request_sanitize"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  rules {
    action     = "execute"
    enabled    = true
    expression = "true"
    id         = "0789dc4343054d1e981f8c44bedc6fbd"
    ref        = "0789dc4343054d1e981f8c44bedc6fbd"
    version    = "1"
    action_parameters {
      overrides {
        rules {
          enabled = true
          id      = "78723a9e0c7c4c6dbec5684cb766231d"
        }
      }
      id      = "70339d97bdb34195bbf054b1ebe81f76"
      version = "latest"
    }
  }
}

resource "cloudflare_ruleset" "terraform_managed_resource" {
  kind    = "zone"
  name    = "default"
  phase   = "http_ratelimit"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  rules {
    action      = "block"
    description = "fwewe"
    enabled     = false
    expression  = "(http.cookie eq \"namwe=value\")"
    id          = "549e64153ff14d2cb5a5ef88c1f5bdbc"
    ref         = "549e64153ff14d2cb5a5ef88c1f5bdbc"
    version     = "1"
    ratelimit {
      characteristics     = ["ip.src", "cf.colo.id"]
      mitigation_timeout  = 30
      period              = 60
      requests_per_period = 100
    }
  }
}

resource "cloudflare_ruleset" "terraform_managed_resource" {
  kind    = "zone"
  name    = "zone"
  phase   = "ddos_l7"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  rules {
    action      = "execute"
    description = "zone"
    enabled     = true
    expression  = "true"
    id          = "c6893ad10fb344e9b8be3c0c3575adc9"
    ref         = "c6893ad10fb344e9b8be3c0c3575adc9"
    version     = "1"
    action_parameters {
      id      = "4d21379b4f9f4bb088e0729962c8b3cf"
      version = "latest"
    }
  }
}

resource "cloudflare_ruleset" "terraform_managed_resource" {
  kind    = "zone"
  name    = "zone"
  phase   = "http_request_firewall_managed"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  rules {
    action      = "execute"
    description = "zone"
    enabled     = false
    expression  = "(http.cookie eq \"jb_testing=true\")"
    id          = "6c41b21b852a4c89a0eaaf5a7ac560a8"
    ref         = "6c41b21b852a4c89a0eaaf5a7ac560a8"
    version     = "2"
    action_parameters {
      overrides {
        rules {
          action = "block"
          id     = "5de7edfa648c4d6891dc3e7f84534ffa"
        }
        rules {
          action = "block"
          id     = "d52aa57408a144afa35e0fd96e3897dc"
        }
        rules {
          action = "block"
          id     = "7994335d116849f7a0ab6b771d1d0db7"
        }
        rules {
          action = "block"
          id     = "20e34d3164a340dbb5c5d29203ccff90"
        }
        rules {
          action = "block"
          id     = "8d9f209f35df412ba4bafe5156335ab1"
        }
        rules {
          action = "block"
          id     = "8840c3fa2c7947f6b10176ceb8f65558"
        }
        rules {
          action = "block"
          id     = "48e06376fc6347c0bf08b8ccf82d008b"
        }
        rules {
          action = "block"
          id     = "8ea0937695984040b528c80a4e6df495"
        }
        rules {
          action = "block"
          id     = "b777ce009bb346b39be4886055a71165"
        }
        rules {
          action = "block"
          id     = "cb5b6de178d3488d8649da8608b7b3a2"
        }
        rules {
          action = "block"
          id     = "390b6273c8dc4366b36e52fc6f35c356"
        }
        rules {
          action = "block"
          id     = "8ac6964456494da6b098a93c35f86fc9"
        }
        rules {
          action = "block"
          id     = "5ac122b3972c4247a247f3271045f374"
        }
        rules {
          action = "block"
          id     = "b1efd337665d49f5950f892971120c4b"
        }
        rules {
          action = "block"
          id     = "34158d546873469a8f8ccee19139627b"
        }
      }
      id      = "efb7b8c949ac4650a09736fc376e9aee"
      version = "latest"
    }
  }
}

resource "cloudflare_ruleset" "terraform_managed_resource" {
  kind    = "zone"
  name    = "default"
  phase   = "http_request_late_transform"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  rules {
    action      = "rewrite"
    description = "test transform"
    enabled     = true
    expression  = "(http.request.uri.path eq \"example.com\")"
    id          = "e5b61605d6cf4ce08f729c17d42d76ef"
    ref         = "e5b61605d6cf4ce08f729c17d42d76ef"
    version     = "1"
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
  }
  rules {
    action      = "rewrite"
    description = "test transform set"
    enabled     = true
    expression  = "(http.request.uri.path eq \"example.com\")"
    id          = "8ec764cf386940c89dd83dbab7bb4c16"
    ref         = "8ec764cf386940c89dd83dbab7bb4c16"
    version     = "1"
    action_parameters {
      headers {
        expression = "(ip.geoip.continent eq \"T1\")"
        name       = "example-http-static-header-1"
        operation  = "set"
        value      = "my-http-header-1"
      }
    }
  }
  rules {
    action      = "rewrite"
    description = "test uri rewrite set"
    enabled     = false
    expression  = "(http.request.uri.path eq \"pumpkin.com\")"
    id          = "d0f1b4fdb4234adf9c6de9b614424836"
    ref         = "d0f1b4fdb4234adf9c6de9b614424836"
    version     = "1"
    action_parameters {
      uri {
        path {
          value = "/spaceship"
        }
      }
    }
  }
}

resource "cloudflare_ruleset" "terraform_managed_resource" {
  description = "Some ruleset"
  kind        = "zone"
  name        = "default"
  phase       = "http_request_late_transform"
  zone_id     = "0da42c8d2132a9ddaf714f9e7c920711"
}
