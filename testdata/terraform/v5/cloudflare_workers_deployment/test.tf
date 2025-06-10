resource "cloudflare_workers_deployment" "terraform_managed_resource_0" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  script_name = "script_2"
  strategy    = "percentage"
  annotations = {
    "workers/message"      = "Automatic deployment on upload."
    "workers/triggered_by" = "upload"
  }
  versions = [{
    percentage = 100
    version_id = "a81e19c8-2d06-4495-b4dd-7fcbf7729a86"
  }]
}

resource "cloudflare_workers_deployment" "terraform_managed_resource_1" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  script_name = "script_2"
  strategy    = "percentage"
  annotations = {
    "workers/message"      = "Automatic deployment on upload."
    "workers/triggered_by" = "upload"
  }
  versions = [{
    percentage = 100
    version_id = "a81e19c8-2d06-4495-b4dd-7fcbf7729a86"
  }]
}

