resource "cloudflare_teams_list" "terraform_managed_resource" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  description = "we like domains here"
  items       = ["example.com"]
  name        = "STUFF TO DO WITH DOMAINS"
  type        = "DOMAIN"
}
