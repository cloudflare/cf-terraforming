resource "cloudflare_load_balancer" "terraform_managed_resource" {
  default_pool_ids = ["17b5962d775c646f3f9725cbc7a53df4", "9290f38c5d07c2e2f4df57b1f61d4196", "00920f38ce07c2e2f4df50b1f61d4194"]
  description      = "Load Balancer for www.example.com"
  enabled          = true
  fallback_pool_id = "17b5962d775c646f3f9725cbc7a53df4"
  name             = "www.example.com"
  proxied          = true
  session_affinity = "cookie"
  session_affinity_attributes = {
    drain_duration = 100
    samesite       = "Auto"
    secure         = "Auto"
  }
  session_affinity_ttl = 5000
  steering_policy      = "dynamic_latency"
  ttl                  = 30
  zone_id              = "0da42c8d2132a9ddaf714f9e7c920711"
  pop_pools {
    pool_ids = ""
    pop      = ""
  }
  region_pools {
    pool_ids = ""
    region   = ""
  }
}
