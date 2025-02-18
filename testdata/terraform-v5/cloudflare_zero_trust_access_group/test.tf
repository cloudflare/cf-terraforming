resource "cloudflare_zero_trust_access_group" "example_zero_trust_access_group" {
  include = [{
    group = {
      id = "aa0a4aab-672b-4bdb-bc33-a59f1130a11f"
    }
  }]
  name = "Allow devs"
  zone_id = "zone_id"
  exclude = [{
    group = {
      id = "aa0a4aab-672b-4bdb-bc33-a59f1130a11f"
    }
  }]
  is_default = true
  require = [{
    group = {
      id = "aa0a4aab-672b-4bdb-bc33-a59f1130a11f"
    }
  }]
}
