resource "cloudflare_ai_gateway" "terraform_managed_resource" {
  account_id                 = "f037e56e89293a057740de681ac9abbe"
  authentication             = false
  cache_invalidate_on_update = false
  cache_ttl                  = 0
  collect_logs               = true
  log_management             = 10000000
  log_management_strategy    = "DELETE_OLDEST"
  logpush                    = false
  rate_limiting_interval     = 0
  rate_limiting_limit        = 0
  rate_limiting_technique    = "fixed"
  workers_ai_billing_mode    = "postpaid"
  zdr                        = false
}
