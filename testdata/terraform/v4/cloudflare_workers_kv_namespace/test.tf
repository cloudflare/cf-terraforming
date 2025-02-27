resource "cloudflare_workers_kv_namespace" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  title      = "example-kv-namespace"
}
