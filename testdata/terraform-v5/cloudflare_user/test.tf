resource "cloudflare_user" "example_user" {
  country = "US"
  first_name = "John"
  last_name = "Appleseed"
  telephone = "+1 123-123-1234"
  zipcode = "12345"
}
