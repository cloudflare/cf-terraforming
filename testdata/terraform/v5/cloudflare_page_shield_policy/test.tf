resource "cloudflare_page_shield_policy" "terraform_managed_resource" {
  action      = "log"
  description = "test-policy"
  enabled     = true
  expression  = "(ip.src.country eq \"CL\")"
  value       = "default-src 'none'"
  zone_id     = "0da42c8d2132a9ddaf714f9e7c920711"
}

