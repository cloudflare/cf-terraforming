resource "cloudflare_user" "terraform_managed_resource" {
  country    = "US"
  first_name = "john"
  last_name  = "doe"
  telephone  = "+1234567890"
  zipcode    = "1234"
}
