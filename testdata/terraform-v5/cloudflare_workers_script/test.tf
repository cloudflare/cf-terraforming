resource "cloudflare_workers_script" "example_workers_script" {
  account_id = "023e105f4ecef8ad9ca31a8372d0c353"
  script_name = "this-is_my_script-01"
  metadata = {
    assets = {
      config = {
        html_handling = "auto-trailing-slash"
        not_found_handling = "none"
        run_worker_first = false
        serve_directly = true
      }
      jwt = "jwt"
    }
    bindings = [{
      name = "MY_ENV_VAR"
      type = "plain_text"
    }]
    body_part = "worker.js"
    compatibility_date = "2021-01-01"
    compatibility_flags = ["nodejs_compat"]
    keep_assets = false
    keep_bindings = ["string"]
    logpush = false
    main_module = "worker.js"
    migrations = {
      deleted_classes = ["string"]
      new_classes = ["string"]
      new_sqlite_classes = ["string"]
      new_tag = "v2"
      old_tag = "v1"
      renamed_classes = [{
        from = "from"
        to = "to"
      }]
      transferred_classes = [{
        from = "from"
        from_script = "from_script"
        to = "to"
      }]
    }
    observability = {
      enabled = true
      head_sampling_rate = 0.1
    }
    placement = {
      mode = "smart"
    }
    tags = ["string"]
    tail_consumers = [{
      service = "my-log-consumer"
      environment = "production"
      namespace = "my-namespace"
    }]
    usage_model = "standard"
  }
}
