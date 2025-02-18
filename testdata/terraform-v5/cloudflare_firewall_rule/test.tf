resource "cloudflare_firewall_rule" "example_firewall_rule" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  action = {
    mode = "simulate"
    response = {
      body = "<error>This request has been rate-limited.</error>"
      content_type = "text/xml"
    }
    timeout = 86400
  }
  filter = {
    description = "Restrict access from these browsers on this address range."
    expression = "(http.request.uri.path ~ \".*wp-login.php\" or http.request.uri.path ~ \".*xmlrpc.php\") and ip.addr ne 172.16.22.155"
    paused = false
    ref = "FIL-100"
  }
}
