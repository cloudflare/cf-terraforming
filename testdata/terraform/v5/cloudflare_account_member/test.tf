resource "cloudflare_account_member" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  email      = "ggalow@cloudflare.com"
  roles      = ["33666b9c79b9a5273fc7344ff42f953d"]
  status     = "accepted"
}