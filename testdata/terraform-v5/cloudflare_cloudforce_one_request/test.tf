resource "cloudflare_cloudforce_one_request" "example_cloudforce_one_request" {
  account_identifier = "023e105f4ecef8ad9ca31a8372d0c353"
  content = "What regions were most effected by the recent DoS?"
  priority = "routine"
  request_type = "Victomology"
  summary = "DoS attack"
  tlp = "clear"
}
