resource "cloudflare_pages_project" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  name = "test_project"
  production_branch = "main"
  build_config {
    build_command = "npm run build"
    destination_dir = "build"
    root_dir = "/"
    web_analytics_tag = "0ee1d926cd60d2618a108d4232a75b73"
    web_analytics_token = "c05bb382259183db3a0a822b64c11459"
  }
  source {
    type = "github"
    config {
      owner = "cloudflare"
      repo_name = "pages-test"
      production_branch = "main"
      pr_comments_enabled = true
      deployments_enabled = true
      production_deployment_enabled = true
      preview_deployment_setting = "custom"
      preview_branch_includes = [ "release/*", "production", "main"]
      preview_branch_excludes = ["dependabot/*", "dev", "*/ignore"]
    }
  }
  deployment_configs {
    preview {
      environment_variables = {
        ENVIRONMENT = "preview"
        BUILD_VERSION = "1.2"
      }
      kv_namespaces = {
        KV_BINDING = "5eb63bbbe01eeed093cb22bb8f5acdc3"
      }
      durable_object_namespaces = {
        DO_BINDING = "5eb63bbbe01eeed093cb22bb8f5acdc3"
      }
      r2_buckets = {
        R2_BINDING = "some-bucket"
      }
      d1_databases = {
        D1_BINDING = "a94509c6-0757-43f3-b053-474b0ab10935"
      }
      compatibility_date = "2022-08-16"
      compatibility_flags = ["preview_flag"]
    }
    production {
      environment_variables = {
        ENVIRONMENT = "production"
        BUILD_VERSION = "1.2"
      }
      kv_namespaces = {
        KV_BINDING_1 = "5eb63bbbe01eeed093cb22bb8f5acdc3"
        KV_BINDING_2 = "3cdca5f8bb22bc390deee10ebbb36be5"
      }
      durable_object_namespaces = {
        DO_BINDING_1 = "5eb63bbbe01eeed093cb22bb8f5acdc3"
        DO_BINDING_2 = "3cdca5f8bb22bc390deee10ebbb36be5"
      }
      r2_buckets = {
        R2_BINDING_1 = "some-bucket"
        R2_BINDING_2 = "other-bucket"
      }
      d1_databases = {
        D1_BINDING_1 = "445e2955-951a-4358-a35b-a4d0c813f63"
        D1_BINDING_2 = "a399414b-c697-409a-a688-377db6433cd9"
      }
      compatibility_date = "2022-08-15"
      compatibility_flags = ["production_flag", "second flag"]
    }
  }
}