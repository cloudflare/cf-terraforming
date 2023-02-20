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
      version = "latest"
    }
    description  = "zone"
    enabled      = false
    expression   = "(http.cookie eq \"jb_testing=true\")"
    last_updated = "2021-09-03T06:42:41.341405Z"
  }
}
