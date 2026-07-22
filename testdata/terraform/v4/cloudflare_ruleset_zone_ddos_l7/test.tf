resource "cloudflare_ruleset" "terraform_managed_resource" {
  kind    = "zone"
  name    = "zone"
  phase   = "ddos_l7"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
  rules {
    action = "execute"
    action_parameters {
      id = "4d21379b4f9f4bb088e0729962c8b3cf"
      overrides {
        sensitivity_level = "medium"
      }
    }
    description = "zone"
    enabled     = true
    expression  = "true"
  }
}
