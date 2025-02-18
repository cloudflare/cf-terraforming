resource "cloudflare_queue_consumer" "example_queue_consumer" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  queue_id = "023e105f4ecef8ad9ca31a8372d0c353"
  dead_letter_queue = "example-queue"
  script_name = "my-consumer-worker"
  settings = {
    batch_size = 50
    max_concurrency = 10
    max_retries = 3
    max_wait_time_ms = 5000
    retry_delay = 10
  }
  type = "worker"
}
