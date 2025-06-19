resource "cloudflare_list_item" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  comment    = "okhejrsmza"
  list_id    = "6cafa626bdb6453fac7a9be3aacf73ca"
  redirect = {
    include_subdomains    = false
    preserve_path_suffix  = false
    preserve_query_string = false
    source_url            = "example.com/"
    status_code           = 301
    subpath_matching      = false
    target_url            = "https://example1.com"
  }
}

