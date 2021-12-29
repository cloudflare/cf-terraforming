resource "cloudflare_worker_route" "terraform_managed_resource" {
  pattern = "example.com/*"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
}
