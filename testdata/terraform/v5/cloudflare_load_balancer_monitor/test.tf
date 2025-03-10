resource "cloudflare_load_balancer_monitor" "terraform_managed_resource_0" {
  account_id       = "f037e56e89293a057740de681ac9abbe"
  allow_insecure   = false
  description      = "this is a very weird load balancer"
  expected_body    = "dead"
  expected_codes   = "5xx"
  follow_redirects = false
  header = {
    Header = ["Host"]
    Values = ["terraform.cfapi.net"]
  }
  interval = 60
  method   = "HEAD"
  path     = "/custom"
  port     = 8080
  retries  = 5
  timeout  = 9
  type     = "http"
}

resource "cloudflare_load_balancer_monitor" "terraform_managed_resource_1" {
  account_id       = "f037e56e89293a057740de681ac9abbe"
  allow_insecure   = false
  expected_body    = "alive"
  expected_codes   = "2xx"
  follow_redirects = false
  header           = {}
  interval         = 60
  method           = "GET"
  path             = "/"
  retries          = 2
  timeout          = 5
  type             = "http"
}

resource "cloudflare_load_balancer_monitor" "terraform_managed_resource_2" {
  account_id       = "f037e56e89293a057740de681ac9abbe"
  allow_insecure   = false
  expected_body    = "alive"
  expected_codes   = "2xx"
  follow_redirects = false
  header           = {}
  interval         = 60
  method           = "GET"
  path             = "/"
  retries          = 2
  timeout          = 5
  type             = "http"
}

resource "cloudflare_load_balancer_monitor" "terraform_managed_resource_3" {
  account_id       = "f037e56e89293a057740de681ac9abbe"
  allow_insecure   = false
  expected_body    = "alive"
  expected_codes   = "2xx"
  follow_redirects = false
  header           = {}
  interval         = 60
  method           = "GET"
  path             = "/"
  retries          = 2
  timeout          = 5
  type             = "http"
}

resource "cloudflare_load_balancer_monitor" "terraform_managed_resource_4" {
  account_id       = "f037e56e89293a057740de681ac9abbe"
  allow_insecure   = false
  expected_body    = "alive"
  expected_codes   = "2xx"
  follow_redirects = false
  header           = {}
  interval         = 60
  method           = "GET"
  path             = "/"
  retries          = 2
  timeout          = 5
  type             = "http"
}

