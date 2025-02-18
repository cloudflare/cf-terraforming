resource "cloudflare_list" "example_list" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  kind = "ip"
  name = "list1"
  description = "This is a note"
}
