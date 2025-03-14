resource "cloudflare_logpush_job" "terraform_managed_resource_0" {
  account_id       = "f037e56e89293a057740de681ac9abbe"
  dataset          = "workers_trace_events"
  destination_conf = "r2://terraform-acctest/date={DATE}?account-id=f037e56e89293a057740de681ac9abbe&access-key-id=0c6710b6f5a77c3b4735f49616694cf2&secret-access-key=7bb918ea5bc7c68729c597a3444e1095b3bbbcaacdce8288f1f243bfd822f337"
  enabled          = true
  frequency        = "high"
  logpull_options  = "fields=Event,EventTimestampMs,Outcome,Exceptions,Logs,ScriptName"
  name             = "fmvpbbpnkb"
}

resource "cloudflare_logpush_job" "terraform_managed_resource_1" {
  account_id       = "f037e56e89293a057740de681ac9abbe"
  dataset          = "workers_trace_events"
  destination_conf = "r2://terraform-acctest/date={DATE}?account-id=f037e56e89293a057740de681ac9abbe&access-key-id=0c6710b6f5a77c3b4735f49616694cf2&secret-access-key=7bb918ea5bc7c68729c597a3444e1095b3bbbcaacdce8288f1f243bfd822f337"
  enabled          = true
  frequency        = "high"
  logpull_options  = "fields=Event,EventTimestampMs,Outcome,Exceptions,Logs,ScriptName"
  name             = "httgotkhpj"
}

