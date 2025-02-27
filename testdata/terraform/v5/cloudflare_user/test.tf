resource "cloudflare_user" "terraform_managed_resource" {
  country = "US"
  first_name = "John"
  last_name = "Appleseed"
  telephone = "+1 123-123-1234"
  zipcode = "12345"
}
