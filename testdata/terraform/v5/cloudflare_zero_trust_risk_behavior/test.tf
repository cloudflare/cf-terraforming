resource "cloudflare_zero_trust_risk_behavior" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  behaviors = {
    high_dlp = {
      description = "User has triggered an active DLP profile in a Gateway policy fifteen times or more within one minute."
      enabled     = true
      name        = "High Number of DLP Policies Triggered"
      risk_level  = "medium"
    }
    imp_travel = {
      description = "A user had a successful Access application log in from two locations that they could not have traveled to in that period of time."
      enabled     = true
      name        = "Impossible Travel"
      risk_level  = "high"
    }
    sentinel_one = {
      description = "User is signed in on a device where Sentinel One EDR detects an active threat or infection."
      enabled     = false
      name        = "SentinelOne Infection Detected"
      risk_level  = "medium"
    }
  }
}

