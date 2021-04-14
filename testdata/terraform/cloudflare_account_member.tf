resource "cloudflare_account_member" "terraform_managed_resource" {
  email_address = "user@example.com"
  role_ids = [ "3536bcfad5faccb999b47003c79917fb" ]
}
