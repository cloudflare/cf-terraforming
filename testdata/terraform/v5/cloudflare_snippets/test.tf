resource "cloudflare_snippets" "terraform_managed_resource" {
  files        = []
  snippet_name = "example_snippet"
  zone_id      = "0da42c8d2132a9ddaf714f9e7c920711"
  metadata = {
    main_module = "main.js"
  }
}

