resource "cloudflare_load_balancer" "terraform_managed_resource" {
  default_pool_ids     = ["17b5962d775c646f3f9725cbc7a53df4", "9290f38c5d07c2e2f4df57b1f61d4196", "00920f38ce07c2e2f4df50b1f61d4194"]
  description          = "Load Balancer for www.example.com"
  enabled              = true
  fallback_pool_id     = "17b5962d775c646f3f9725cbc7a53df4"
  name                 = "www.example.com"
  proxied              = false
  session_affinity     = "cookie"
  session_affinity_ttl = 5000
  steering_policy      = "dynamic_latency"
  ttl                  = 30
  zone_id              = "0da42c8d2132a9ddaf714f9e7c920711"
  country_pools {
    country  = "AU"
    pool_ids = ["de90f38ced07c2e2f4df50b1f61d4194", "9290f38c5d07c2e2f4df57b1f61d4196"]
  }
  pop_pools {
    pool_ids = ["00920f38ce07c2e2f4df50b1f61d4194"]
    pop      = "SJC"
  }
  random_steering {
    default_weight = 1
    pool_weights = {
      "2c3f886957b4112bfaca8b12d87ce8c1" = 0
    }
  }
  region_pools {
    pool_ids = ["00920f38ce07c2e2f4df50b1f61d4194"]
    region   = "ENAM"
  }
  session_affinity_attributes {
    drain_duration = 100
    samesite       = "Auto"
    secure         = "Auto"
  }
}
