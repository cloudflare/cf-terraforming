resource "cloudflare_workers_cron_trigger" "terraform_managed_resource" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  script_name = "script_2"
  schedules = [{
    cron = "*/30 * * * *"
  }]
}

