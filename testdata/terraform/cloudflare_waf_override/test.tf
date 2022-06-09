resource "cloudflare_waf_override" "terraform_managed_resource" {
  description = "Enable Cloudflare Magento ruleset for shop.example.com"
  groups = {
    ea8687e59929c1fd05ba97574ad43f77 = "default"
  }
  paused   = false
  priority = 1
  rewrite_action = {
    challenge = "block"
    default   = "block"
    simulate  = "disable"
  }
  rules = {
    100015 = "disable"
  }
  urls    = ["shop.example.com/*"]
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
}
