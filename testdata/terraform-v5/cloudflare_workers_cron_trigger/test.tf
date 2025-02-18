resource "cloudflare_workers_cron_trigger" "example_workers_cron_trigger" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  script_name = "this-is_my_script-01"
  body = [{
    cron = "*/30 * * * *"
  }]
}
