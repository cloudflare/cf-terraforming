resource "cloudflare_workers_custom_domain" "example_workers_custom_domain" {
  account_id = "9a7806061c88ada191ed06f989cc3dac"
  environment = "production"
  hostname = "foo.example.com"
  service = "foo"
  zone_id = "593c9c94de529bbbfaac7c53ced0447d"
}
