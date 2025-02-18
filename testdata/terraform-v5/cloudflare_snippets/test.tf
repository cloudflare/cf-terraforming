resource "cloudflare_snippets" "example_snippets" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  snippet_name = "snippet_name_01"
  files = "export { async function fetch(request, env) {return new Response(\'some_response\') } }"
  metadata = {
    main_module = "main.js"
  }
}
