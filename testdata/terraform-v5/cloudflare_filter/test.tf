resource "cloudflare_filter" "example_filter" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  expression = "(http.request.uri.path ~ \".*wp-login.php\" or http.request.uri.path ~ \".*xmlrpc.php\") and ip.addr ne 172.16.22.155"
}
