resource "cloudflare_load_balancer" "terraform_managed_resource_0" {
  default_pools    = ["c36b8a3066b335b2af7e940f2588805d"]
  enabled          = true
  fallback_pool    = "c36b8a3066b335b2af7e940f2588805d"
  name             = "tf-testacc-lb-fwwjdnedoi.terraform.cfapi.net"
  networks         = ["cloudflare"]
  pop_pools        = {}
  proxied          = false
  region_pools     = {}
  session_affinity = "none"
  steering_policy  = "off"
  ttl              = 30
  zone_id          = "0da42c8d2132a9ddaf714f9e7c920711"
  adaptive_routing = {
    failover_across_pools = false
  }
  location_strategy = {
    mode       = "pop"
    prefer_ecs = "proximity"
  }
  random_steering = {
    default_weight = 1
  }
  session_affinity_attributes = {
    drain_duration         = 0
    samesite               = "Auto"
    secure                 = "Auto"
    zero_downtime_failover = "none"
  }
}

resource "cloudflare_load_balancer" "terraform_managed_resource_1" {
  default_pools    = ["0ce4832a7181e0c3e2936e2c34a4687f"]
  description      = "rules lb"
  enabled          = true
  fallback_pool    = "0ce4832a7181e0c3e2936e2c34a4687f"
  name             = "tf-testacc-lb-sidcrfxrak.terraform.cfapi.net"
  networks         = ["cloudflare"]
  pop_pools        = {}
  proxied          = false
  region_pools     = {}
  session_affinity = "none"
  steering_policy  = "off"
  ttl              = 30
  zone_id          = "0da42c8d2132a9ddaf714f9e7c920711"
  adaptive_routing = {
    failover_across_pools = false
  }
  location_strategy = {
    mode       = "pop"
    prefer_ecs = "proximity"
  }
  random_steering = {
    default_weight = 1
  }
  rules = [{
    condition = "dns.qry.type == 28"
    disabled  = false
    name      = "test rule 1"
    overrides = {
      adaptive_routing = {
        failover_across_pools = true
      }
      location_strategy = {
        mode       = "resolver_ip"
        prefer_ecs = "always"
      }
      random_steering = {
        default_weight = 0.2
        pool_weights = {
          c29c1dc121903fbea9f0c92e83a1b1e2 = 0.4
        }
      }
      session_affinity_attributes = {
        require_all_headers    = false
        samesite               = "Auto"
        secure                 = "Auto"
        zero_downtime_failover = "sticky"
      }
      steering_policy = "geo"
    }
    priority = 0
    }, {
    condition = "dns.qry.type == 28"
    disabled  = false
    fixed_response = {
      content_type = "html"
      location     = "www.example.com"
      message_body = "hello"
      status_code  = 200
    }
    name       = "test rule 2"
    overrides  = {}
    priority   = 10
    terminates = true
    }, {
    condition = "dns.qry.type == 28"
    disabled  = false
    name      = "test rule 3"
    overrides = {
      region_pools = {
        ENAM = ["0ce4832a7181e0c3e2936e2c34a4687f"]
      }
    }
    priority = 20
  }]
  session_affinity_attributes = {
    drain_duration         = 0
    samesite               = "Auto"
    secure                 = "Auto"
    zero_downtime_failover = "none"
  }
}

