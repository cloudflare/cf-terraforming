resource "cloudflare_workers_deployment" "example_workers_deployment" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  script_name = "this-is_my_script-01"
  strategy = "percentage"
  versions = [{
    percentage = 100
    version_id = "bcf48806-b317-4351-9ee7-36e7d557d4de"
  }]
  annotations = {
    workers_message = "Deploy bug fix."
  }
}
