resource "cloudflare_dns_firewall" "terraform_managed_resource" {
  account_id             = "f037e56e89293a057740de681ac9abbe"
  deprecate_any_requests = true
  ecs_fallback           = false
  maximum_cache_ttl      = 900
  minimum_cache_ttl      = 60
  name                   = "ygpvauebcd"
  ratelimit              = 1000
  retries                = 2
  upstream_ips           = ["1.2.3.4"]
  attack_mitigation      = {}
}
