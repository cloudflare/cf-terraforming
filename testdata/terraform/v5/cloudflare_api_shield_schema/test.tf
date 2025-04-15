resource "cloudflare_api_shield_schema" "terraform_managed_resource_0" {
  file               = "{\"info\":{\"title\":\"Example\",\"version\":\"0.1.0\"},\"openapi\":\"3.0.3\",\"paths\":{\"/\":{}},\"servers\":[{\"url\":\"api.example.com\"}]}"
  kind               = "openapi_v3"
  name               = "example_schema.json"
  schema_id          = "59f6e0a9-7d8d-446f-b4c8-fb9c2c1abae8"
  validation_enabled = true
  zone_id            = "0da42c8d2132a9ddaf714f9e7c920711"
}

resource "cloudflare_api_shield_schema" "terraform_managed_resource_1" {
  file               = "{\"info\":{\"title\":\"Example\",\"version\":\"0.1.0\"},\"openapi\":\"3.0.3\",\"paths\":{\"/\":{}},\"servers\":[{\"url\":\"api.example.com\"}]}"
  kind               = "openapi_v3"
  name               = "example_schema.json"
  schema_id          = "fef87c1a-6ff7-4d3a-aeee-fe6cf9ca948a"
  validation_enabled = true
  zone_id            = "0da42c8d2132a9ddaf714f9e7c920711"
}

