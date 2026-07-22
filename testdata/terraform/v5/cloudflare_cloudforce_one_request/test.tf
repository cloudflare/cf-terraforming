resource "cloudflare_cloudforce_one_request" "terraform_managed_resource" {
  account_identifier = "023e105f4ecef8ad9ca31a8372d0c353"
  content = "What regions were most effected by the recent DoS?"
  priority = "routine"
  request_type = "Victomology"
  summary = "DoS attack"
  tlp = "clear"
}
