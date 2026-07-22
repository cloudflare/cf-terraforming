resource "cloudflare_ruleset" "terraform_managed_resource" {
  kind    = "zone"
  name    = "zone"
  phase   = "http_request_firewall_managed"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  rules {
    action = "execute"
    action_parameters {
      id = "efb7b8c949ac4650a09736fc376e9aee"
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
    }
    description = "zone"
    enabled     = false
    expression  = "(http.cookie eq \"jb_testing=true\")"
  }
  rules {
    action = "skip"
    action_parameters {
      rules = {
        "4814384a9e5d4991b9815dcfc25d2f1f" = "37da7855d2f94f69865365d894a556a4,6afe6795ee6a48d6a1dfe59255395a78,5a6f5a57cde8428ab0668ce17cdec0c8,5e4903d6afa841c9b88b96203297003f,2380cd409b604c2a9273042f3eb29c4e,f5aebedc99a14c8d9e8cfa2ce5f94216,edf8c37cc81747d382690b3c77e82ce4,1129dfb383bb42e48466488cf3b37cb1"
      }
    }
    description = "Bypass managed OWSAP SQL Injection rules for /api/v1/identity"
    enabled     = true
    expression  = "(http.request.method eq \"POST\" and http.request.uri.path eq \"/api/v1/identity\")"
    logging {
      enabled = true
    }
  }
}
