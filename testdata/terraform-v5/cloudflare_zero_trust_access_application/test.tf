resource "cloudflare_zero_trust_access_application" "example_zero_trust_access_application" {
  domain = "test.example.com/admin"
  type = "self_hosted"
  zone_id = "zone_id"
  allow_authenticate_via_warp = true
  allowed_idps = ["699d98642c564d2e855e9661899b7252"]
  app_launcher_visible = true
  auto_redirect_to_identity = true
  cors_headers = {
    allow_all_headers = true
    allow_all_methods = true
    allow_all_origins = true
    allow_credentials = true
    allowed_headers = ["string"]
    allowed_methods = ["GET"]
    allowed_origins = ["https://example.com"]
    max_age = -1
  }
  custom_deny_message = "custom_deny_message"
  custom_deny_url = "custom_deny_url"
  custom_non_identity_deny_url = "custom_non_identity_deny_url"
  custom_pages = ["699d98642c564d2e855e9661899b7252"]
  destinations = [{
    type = "public"
    uri = "test.example.com/admin"
  }, {
    type = "public"
    uri = "test.anotherexample.com/staff"
  }, {
    cidr = "10.5.0.0/24"
    hostname = "hostname"
    l4_protocol = "tcp"
    port_range = "80-90"
    type = "private"
    vnet_id = "vnet_id"
  }, {
    cidr = "10.5.0.3/32"
    hostname = "hostname"
    l4_protocol = "tcp"
    port_range = "80"
    type = "private"
    vnet_id = "vnet_id"
  }, {
    cidr = "cidr"
    hostname = "hostname"
    l4_protocol = "tcp"
    port_range = "port_range"
    type = "private"
    vnet_id = "vnet_id"
  }]
  enable_binding_cookie = true
  http_only_cookie_attribute = true
  logo_url = "https://www.cloudflare.com/img/logo-web-badges/cf-logo-on-white-bg.svg"
  name = "Admin Site"
  options_preflight_bypass = true
  path_cookie_attribute = true
  policies = [{
    id = "f174e90a-fafe-4643-bbbc-4a0ed4fc8415"
    precedence = 0
  }]
  same_site_cookie_attribute = "strict"
  scim_config = {
    idp_uid = "idp_uid"
    remote_uri = "remote_uri"
    authentication = {
      password = "password"
      scheme = "httpbasic"
      user = "user"
    }
    deactivate_on_delete = true
    enabled = true
    mappings = [{
      schema = "urn:ietf:params:scim:schemas:core:2.0:User"
      enabled = true
      filter = "title pr or userType eq \"Intern\""
      operations = {
        create = true
        delete = true
        update = true
      }
      strictness = "strict"
      transform_jsonata = "$merge([$, {\'userName\': $substringBefore($.userName, \'@\') & \'+test@\' & $substringAfter($.userName, \'@\')}])"
    }]
  }
  self_hosted_domains = ["test.example.com/admin", "test.anotherexample.com/staff"]
  service_auth_401_redirect = true
  session_duration = "24h"
  skip_interstitial = true
  tags = ["engineers"]
}
