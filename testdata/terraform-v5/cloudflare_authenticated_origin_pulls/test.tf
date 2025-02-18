resource "cloudflare_authenticated_origin_pulls" "example_authenticated_origin_pulls" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  config = [{
    cert_id = "2458ce5a-0c35-4c7f-82c7-8e9487d3ff60"
    enabled = true
    hostname = "app.example.com"
  }]
}
