resource "cloudflare_zero_trust_tunnel_cloudflared_config" "example_zero_trust_tunnel_cloudflared_config" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  tunnel_id = "f70ff985-a4ef-4643-bbbc-4a0ed4fc8415"
  config = {
    ingress = [{
      hostname = "tunnel.example.com"
      service = "https://localhost:8001"
      origin_request = {
        access = {
          aud_tag = ["string"]
          team_name = "teamName"
          required = true
        }
        ca_pool = "caPool"
        connect_timeout = 0
        disable_chunked_encoding = true
        http2_origin = true
        http_host_header = "httpHostHeader"
        keep_alive_connections = 0
        keep_alive_timeout = 0
        no_happy_eyeballs = true
        no_tls_verify = true
        origin_server_name = "originServerName"
        proxy_type = "proxyType"
        tcp_keep_alive = 0
        tls_timeout = 0
      }
      path = "subpath"
    }]
    origin_request = {
      access = {
        aud_tag = ["string"]
        team_name = "teamName"
        required = true
      }
      ca_pool = "caPool"
      connect_timeout = 0
      disable_chunked_encoding = true
      http2_origin = true
      http_host_header = "httpHostHeader"
      keep_alive_connections = 0
      keep_alive_timeout = 0
      no_happy_eyeballs = true
      no_tls_verify = true
      origin_server_name = "originServerName"
      proxy_type = "proxyType"
      tcp_keep_alive = 0
      tls_timeout = 0
    }
  }
}
