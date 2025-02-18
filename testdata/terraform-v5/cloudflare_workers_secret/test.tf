resource "cloudflare_workers_secret" "example_workers_secret" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  dispatch_namespace = "my-dispatch-namespace"
  script_name = "this-is_my_script-01"
  name = "MY_SECRET"
  text = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
  type = "secret_text"
}
