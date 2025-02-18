resource "cloudflare_cloud_connector_rules" "example_cloud_connector_rules" {
  zone_id = "023e105f4ecef8ad9ca31a8372d0c353"
  rules = [{
    id = "95c365e17e1b46599cd99e5b231fac4e"
    description = "Rule description"
    enabled = true
    expression = "http.cookie eq \"a=b\""
    parameters = {
      host = "examplebucket.s3.eu-north-1.amazonaws.com"
    }
    cloud_provider = "aws_s3"
  }]
}
