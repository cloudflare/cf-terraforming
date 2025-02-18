resource "cloudflare_zero_trust_device_posture_rule" "example_zero_trust_device_posture_rule" {
  account_id = "699d98642c564d2e855e9661899b7252"
  name = "Admin Serial Numbers"
  type = "file"
  description = "The rule for admin serial numbers"
  expiration = "1h"
  input = {
    operating_system = "windows"
    path = "/bin/cat"
    exists = true
    sha256 = "https://api.us-2.crowdstrike.com"
    thumbprint = "0aabab210bdb998e9cf45da2c9ce352977ab531c681b74cf1e487be1bbe9fe6e"
  }
  match = [{
    platform = "windows"
  }]
  schedule = "1h"
}
