resource "cloudflare_account_member" "terraform_managed_resource" {
  account_id    = "f037e56e89293a057740de681ac9abbe"
  email_address = "user@example.com"
  role_ids      = ["3536bcfad5faccb999b47003c79917fb"]
}
