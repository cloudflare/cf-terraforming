resource "cloudflare_list_item" "example_list_item" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  list_id = "2c0fc9fa937b11eaa1b71c4d701ab86e"
  body = [{
    asn = 5567
    comment = "Private IP address"
    hostname = {
      url_hostname = "example.com"
    }
    ip = "10.0.0.1"
    redirect = {
      source_url = "example.com/arch"
      target_url = "https://archlinux.org/"
      include_subdomains = true
      preserve_path_suffix = true
      preserve_query_string = true
      status_code = 301
      subpath_matching = true
    }
  }]
}
