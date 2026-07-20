resource "cloudflare_ai_search_instance" "terraform_managed_resource" {
  account_id      = "f037e56e89293a057740de681ac9abbe"
  cache           = true
  cache_threshold = "close_enough"
  cache_ttl       = 172800
  chunk           = true
  chunk_overlap   = 10
  chunk_size      = 256
  fusion_method   = "rrf"
  max_num_results = 10
  paused          = false
  reranking       = false
  rewrite_query   = false
  score_threshold = 0.4
  summarization   = false
  sync_interval   = 21600
}
