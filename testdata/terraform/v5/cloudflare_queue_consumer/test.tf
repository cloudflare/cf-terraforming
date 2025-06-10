resource "cloudflare_queue_consumer" "terraform_managed_resource" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  consumer_id = "2f4e3428eaa0472bb6954cf6b7fb932f"
  queue_id    = "2dde6ac405cd457c9ce59dc4bda20c65"
  type        = "worker"
  settings = {
    batch_size       = 50
    max_concurrency  = 10
    max_retries      = 5
    max_wait_time_ms = 5000
    retry_delay      = 10
  }
}

