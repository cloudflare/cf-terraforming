resource "cloudflare_account_member" "example_account_member" {
  account_id = "eb78d65290b24279ba6f44721b3ea3c4"
  email = "user@example.com"
  roles = ["3536bcfad5faccb999b47003c79917fb"]
  status = "accepted"
}
