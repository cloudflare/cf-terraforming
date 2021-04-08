resource "cloudflare_access_rule" "terraform_managed_resource" {
  configuration = {
    target = "ip"
    value = "198.51.100.4"
  }
  mode = "challenge"
  notes = "This rule is on because of an event that occured on date X"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
}
