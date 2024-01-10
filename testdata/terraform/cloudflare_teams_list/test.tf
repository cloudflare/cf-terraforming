resource "cloudflare_teams_list" "terraform_managed_resource" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  name        = "Admin Serial Numbers"
  type        = "SERIAL"
  description = "Serial numbers for all administrators."
  items       = ["8GE8721REF"]
}
