resource "cloudflare_leaked_credential_check_rule" "terraform_managed_resource_0" {
  password = "lookup_json_string(http.request.body.raw, \"password\")"
  username = "lookup_json_string(http.request.body.raw, \"username\")"
  zone_id  = "0da42c8d2132a9ddaf714f9e7c920711"
}

resource "cloudflare_leaked_credential_check_rule" "terraform_managed_resource_1" {
  password = "lookup_json_string(http.request.body.raw, \"pass\")"
  username = "lookup_json_string(http.request.body.raw, \"user\")"
  zone_id  = "0da42c8d2132a9ddaf714f9e7c920711"
}

resource "cloudflare_leaked_credential_check_rule" "terraform_managed_resource_2" {
  password = "lookup_json_string(http.request.body.raw, \"secret\")"
  username = "lookup_json_string(http.request.body.raw, \"id\")"
  zone_id  = "0da42c8d2132a9ddaf714f9e7c920711"
}

