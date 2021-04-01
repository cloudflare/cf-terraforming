resource "cloudflare_load_balancer" "terraform_managed_resource" {
  created_on = "2014-01-01T05:20:00.12345Z"
  description = "Load Balancer for www.example.com"
  enabled = true
  modified_on = "2014-01-01T05:20:00.12345Z"
  name = "www.example.com"
  proxied = true
  session_affinity = "cookie"
  session_affinity_attributes = {
    drain_duration = 100
    samesite = "Auto"
    secure = "Auto"
  }
  session_affinity_ttl = 5000
  steering_policy = "dynamic_latency"
  ttl = 30
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
}
