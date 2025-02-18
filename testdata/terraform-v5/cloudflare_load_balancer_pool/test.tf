resource "cloudflare_load_balancer_pool" "example_load_balancer_pool" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  name = "primary-dc-1"
  origins = [{
    address = "0.0.0.0"
    enabled = true
    header = {
      host = ["example.com"]
    }
    name = "app-server-1"
    virtual_network_id = "a5624d4e-044a-4ff0-b3e1-e2465353d4b4"
    weight = 0.6
  }]
  description = "Primary data center - Provider XYZ"
  enabled = false
  latitude = 0
  load_shedding = {
    default_percent = 0
    default_policy = "random"
    session_percent = 0
    session_policy = "hash"
  }
  longitude = 0
  minimum_origins = 0
  monitor = "monitor"
  notification_email = "someone@example.com,sometwo@example.com"
  notification_filter = {
    origin = {
      disable = true
      healthy = true
    }
    pool = {
      disable = true
      healthy = false
    }
  }
  origin_steering = {
    policy = "random"
  }
}
