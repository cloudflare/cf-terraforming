resource "cloudflare_list" "terraform_managed_resource" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  description = "This is a note"
  kind        = "redirect"
  name        = "redirect_list"
  item {
    comment = "one"
    value {
      redirect {
        include_subdomains    = "enabled"
        preserve_path_suffix  = "disabled"
        preserve_query_string = "enabled"
        source_url            = "example.com/foo"
        status_code           = 301
        subpath_matching      = "enabled"
        target_url            = "https://foo.example.com"
      }
    }
  }
}
