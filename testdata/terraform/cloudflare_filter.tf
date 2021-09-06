resource "cloudflare_filter" "terraform_managed_resource" {
  description = "Restrict access from these browsers on this address range."
  expression  = "(http.request.uri.path ~ \".*wp-login.php\" or http.request.uri.path ~ \".*xmlrpc.php\") and ip.addr ne 172.16.22.155"
  paused      = false
  ref         = "FIL-100"
  zone_id     = "0da42c8d2132a9ddaf714f9e7c920711"
}
