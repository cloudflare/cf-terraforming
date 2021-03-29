resource "cloudflare_byo_ip_prefix" "terraform_managed_resource" {
  advertisement = "on"
  description = "Internal test prefix"
  prefix_id = "9a7806061c88ada191ed06f989cc3dac"
}
