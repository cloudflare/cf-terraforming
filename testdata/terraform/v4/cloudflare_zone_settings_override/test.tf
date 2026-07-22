resource "cloudflare_zone_settings_override" "terraform_managed_resource" {
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  settings {
    always_online             = "on"
    always_use_https          = "off"
    automatic_https_rewrites  = "on"
    brotli                    = "on"
    browser_cache_ttl         = 14400
    browser_check             = "on"
    cache_level               = "aggressive"
    challenge_ttl             = 2700
    ciphers                   = ["ECDHE-RSA-AES128-GCM-SHA256", "AES128-SHA"]
    cname_flattening          = "flatten_at_root"
    development_mode          = "off"
    email_obfuscation         = "on"
    filter_logs_to_cloudflare = "off"
    hotlink_protection        = "off"
    http2                     = "on"
    http3                     = "off"
    ip_geolocation            = "on"
    ipv6                      = "on"
    log_to_cloudflare         = "on"
    max_upload                = 100
    min_tls_version           = "1.2"
    minify {
      css  = "on"
      html = "off"
      js   = "off"
    }
    mirage                      = "off"
    opportunistic_encryption    = "on"
    opportunistic_onion         = "on"
    orange_to_orange            = "off"
    origin_error_page_pass_thru = "off"
    polish                      = "off"
    prefetch_preload            = "off"
    privacy_pass                = "on"
    pseudo_ipv4                 = "off"
    response_buffering          = "off"
    rocket_loader               = "off"
    security_header {
      enabled            = true
      include_subdomains = true
      max_age            = 86400
      nosniff            = true
      preload            = true
    }
    security_level              = "high"
    server_side_exclude         = "on"
    sort_query_string_for_cache = "off"
    ssl                         = "strict"
    tls_1_3                     = "on"
    tls_client_auth             = "off"
    true_client_ip_header       = "off"
    visitor_ip                  = "on"
    waf                         = "on"
    webp                        = "off"
    websockets                  = "on"
    zero_rtt                    = "off"
  }
}
