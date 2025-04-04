resource "cloudflare_load_balancer_pool" "terraform_managed_resource_0" {
  account_id      = "f037e56e89293a057740de681ac9abbe"
  enabled         = true
  minimum_origins = 1
  name            = "pool1"
  origins = [{
    address = "example.com"
    enabled = true
    name    = "example-1"
    weight  = 1
  }]
}

resource "cloudflare_load_balancer_pool" "terraform_managed_resource_1" {
  account_id      = "f037e56e89293a057740de681ac9abbe"
  enabled         = true
  minimum_origins = 1
  name            = "pool2"
  origins = [{
    address = "example.com"
    enabled = true
    name    = "example-2"
    weight  = 1
  }]
}

