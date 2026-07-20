resource "cloudflare_ai_search_instance" "terraform_managed_resource" {
  account_id            = "f037e56e89293a057740de681ac9abbe"
  ai_gateway_id         = "default"
  cache                 = true
  cache_threshold       = "close_enough"
  cache_ttl             = 172800
  chunk                 = true
  chunk_overlap         = 10
  chunk_size            = 1024
  embedding_model       = "@cf/qwen/qwen3-embedding-0.6b"
  fusion_method         = "rrf"
  hybrid_search_enabled = false
  max_num_results       = 10
  paused                = false
  reranking             = false
  rewrite_query         = false
  score_threshold       = 0.4
  summarization         = false
  sync_interval         = 21600
  index_method = {
    keyword = false
    vector  = true
  }
  indexing_options = {
    keyword_tokenizer = "porter"
  }
  retrieval_options = {
    keyword_match_mode = "and"
  }
}
