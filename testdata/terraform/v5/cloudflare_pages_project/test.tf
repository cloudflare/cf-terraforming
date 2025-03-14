resource "cloudflare_pages_project" "terraform_managed_resource_0" {
  account_id        = "f037e56e89293a057740de681ac9abbe"
  name              = "mlfinedniz"
  production_branch = "main"
  build_config = {
    build_command       = ""
    destination_dir     = ""
    root_dir            = ""
    web_analytics_tag   = ""
    web_analytics_token = ""
  }
  deployment_configs = {
    preview = {
      always_use_latest_compatibility_date = false
      build_image_major_version            = 1
      compatibility_flags                  = []
      env_vars                             = null
      fail_open                            = true
      usage_model                          = "standard"
    }
    production = {
      always_use_latest_compatibility_date = false
      build_image_major_version            = 1
      compatibility_flags                  = []
      env_vars                             = null
      fail_open                            = true
      usage_model                          = "standard"
    }
  }
}

resource "cloudflare_pages_project" "terraform_managed_resource_1" {
  account_id        = "f037e56e89293a057740de681ac9abbe"
  name              = "uquivnkfgv"
  production_branch = "main"
  build_config = {
    build_command       = ""
    destination_dir     = ""
    root_dir            = ""
    web_analytics_tag   = ""
    web_analytics_token = ""
  }
  deployment_configs = {
    preview = {
      always_use_latest_compatibility_date = false
      build_image_major_version            = 1
      compatibility_flags                  = []
      env_vars                             = null
      fail_open                            = true
      usage_model                          = "standard"
    }
    production = {
      always_use_latest_compatibility_date = false
      build_image_major_version            = 1
      compatibility_flags                  = []
      env_vars                             = null
      fail_open                            = true
      usage_model                          = "standard"
    }
  }
}

