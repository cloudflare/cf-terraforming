resource "cloudflare_user_agent_blocking_rule" "terraform_managed_resource" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  configuration = {
    target = "ip"
    value = "198.51.100.4"
  }
  mode = "block"
}
