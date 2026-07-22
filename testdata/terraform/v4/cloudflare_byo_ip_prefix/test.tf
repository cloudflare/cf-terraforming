resource "cloudflare_byo_ip_prefix" "terraform_managed_resource" {
  account_id    = "f037e56e89293a057740de681ac9abbe"
  advertisement = "on"
  description   = "Internal test prefix"
  prefix_id     = "9a7806061c88ada191ed06f989cc3dac"
}
