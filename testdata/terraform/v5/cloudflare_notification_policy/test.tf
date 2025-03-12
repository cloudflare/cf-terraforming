resource "cloudflare_notification_policy" "terraform_managed_resource_0" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  alert_type  = "universal_ssl_event_type"
  description = "test description update"
  enabled     = true
  name        = "foo2"
  filters     = {}
  mechanisms = {
    email = [{
      id = "test@example.com"
    }]
  }
}

resource "cloudflare_notification_policy" "terraform_managed_resource_1" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  alert_type  = "zone_aop_custom_certificate_expiration_type"
  description = "This notification is automatically set by Cloudflare"
  enabled     = true
  name        = "Default notification"
  filters     = {}
  mechanisms = {
    email = [{
      id = "test2@example.com"
      }, {
      id = "test3@example.com"
      }, {
      id = "test4@example.com"
      }, {
      id = "test5@example.com"
    }]
  }
}

resource "cloudflare_notification_policy" "terraform_managed_resource_2" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  alert_type  = "billing_usage_alert"
  description = "test description"
  enabled     = true
  name        = "workers usage notification"
  filters = {
    limit   = ["100"]
    product = ["worker_requests"]
  }
  mechanisms = {
    email = [{
      id = "test@example.com"
    }]
  }
}

