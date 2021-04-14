resource "cloudflare_load_balancer" "terraform_managed_resource" {
  description = "Load Balancer for www.example.com"
  enabled = true
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
