resource "cloudflare_workers_kv" "example_workers_kv" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  namespace_id = "0f2ac74b498b48028cb68387c421e279"
  key_name = "My-Key"
  metadata = "{\"someMetadataKey\": \"someMetadataValue\"}"
  value = "Some Value"
}
