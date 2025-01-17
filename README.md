# Cloudflare Terraforming

## Overview

`cf-terraforming` is a command line utility to facilitate terraforming your
existing Cloudflare resources. It does this by using your account credentials to
retrieve your configurations from the [Cloudflare API](https://api.cloudflare.com)
and converting them to Terraform configurations that can be used with the
[Terraform Cloudflare provider](https://registry.terraform.io/providers/cloudflare/cloudflare/latest).

This tool is ideal if you already have Cloudflare resources defined but want to
start managing them via Terraform, and don't want to spend the time to manually
write the Terraform configuration to describe them.

Read the [announcement blog](https://blog.cloudflare.com/cloudflares-partnership-with-hashicorp-and-bootstrapping-terraform-with-cf-terraforming/) for further details on using `cf-terraforming` in your workflow.

> [!NOTE]
> If you would like to export resources compatible with Terraform < 0.12.x,
> you will need to download an older release as this tool no longer supports it.

## Usage

```
Usage:
  cf-terraforming [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  generate    Fetch resources from the Cloudflare API and generate the respective Terraform stanzas
  help        Help about any command
  import      Output `terraform import` compatible commands in order to import resources into state
  version     Print the version number of cf-terraforming

Flags:
  -a, --account string                      Target the provided account ID for the command
  -c, --config string                       Path to config file (default "~/.cf-terraforming.yaml")
  -e, --email string                        API Email address associated with your account
  -h, --help                                help for cf-terraforming
      --hostname string                     Hostname to use to query the API
  -k, --key string                          API Key generated on the 'My Profile' page. See: https://dash.cloudflare.com/profile
      --modern-import-block                 Whether to generate HCL import blocks for generated resources instead of terraform import compatible CLI commands. This is only compatible with Terraform 1.5+
      --provider-registry-hostname string   Hostname to use for provider registry lookups (default "registry.terraform.io")
      --resource-type string                Comma delimitered string of which resource(s) you wish to generate
      --terraform-binary-path string        Path to an existing Terraform binary (otherwise, one will be downloaded)
      --terraform-install-path string       Path to an initialized Terraform working directory (default ".")
  -t, --token string                        API Token
  -v, --verbose                             Specify verbose output (same as setting log level to debug)
  -z, --zone string                         Target the provided zone ID for the command

Use "cf-terraforming [command] --help" for more information about a command.
```

## Authentication

Cloudflare supports two authentication methods to the API:

- API Token - gives access only to resources and permissions specified for that token (recommended)
- API key - gives access to everything your user profile has access to

Both can be retrieved on the [user profile page](https://dash.cloudflare.com/profile/api-tokens).

> [!TIP]
> We recommend that you store your Cloudflare credentials (API key, email, token) as environment
> variables as demonstrated below.

```bash
# if using API Token
export CLOUDFLARE_API_TOKEN='Hzsq3Vub-7Y-hSTlAaLH3Jq_YfTUOCcgf22_Fs-j'

# if using API Key
export CLOUDFLARE_EMAIL='user@example.com'
export CLOUDFLARE_API_KEY='1150bed3f45247b99f7db9696fffa17cbx9'

# specify zone ID
export CLOUDFLARE_ZONE_ID='81b06ss3228f488fh84e5e993c2dc17'

# now call cf-terraforming, e.g.
cf-terraforming generate \
  --resource-type "cloudflare_record" \
  --zone $CLOUDFLARE_ZONE_ID
```

cf-terraforming supports the following environment variables:

- CLOUDFLARE_API_TOKEN - API Token based authentication
- CLOUDFLARE_EMAIL, CLOUDFLARE_API_KEY - API Key based authentication

Alternatively, if using a config file, then specify the inputs using the same
names the `flag` names. Example:

```
cat ~/.cf-terraforming.yaml
email: "email@domain.com"
key: "<key>"
#or
token: "<token>"
```

## Example usage

```bash
cf-terraforming generate \
  --zone $CLOUDFLARE_ZONE_ID \
  --resource-type "cloudflare_record"
```

will contact the Cloudflare API on your behalf and result in a valid Terraform
configuration representing the **resource** you requested:

```hcl
resource "cloudflare_record" "terraform_managed_resource" {
  name = "example.com"
  proxied = false
  ttl = 120
  type = "A"
  value = "198.51.100.4"
  zone_id = "0da42c8d2132a9ddaf714f9e7c920711"
}
```

## Prerequisites

- A Cloudflare account with resources defined (e.g. a few zones, some load
  balancers, spectrum applications, etc)
- A valid Cloudflare API key and sufficient permissions to access the resources
  you are requesting via the API
- An initialised Terraform directory (`terraform init` has run and providers installed). See the [provider documentation](https://registry.terraform.io/providers/cloudflare/cloudflare/latest/docs) if you have not yet setup the Terraform directory.

## Installation

### Homebrew

```bash
brew tap cloudflare/cloudflare
brew install cloudflare/cloudflare/cf-terraforming
```

> [!NOTE]
> If you have installed an older version of `cf-terraforming` via Homebrew,
> you may need to first uninstall `cf-terraforming` and then install it to
> pick up the updated install process and address the signing/notarisation
> issues.

### Go

```bash
go install github.com/cloudflare/cf-terraforming/cmd/cf-terraforming@latest
```

If you use another OS, you will need to download the release directly from
[GitHub Releases](https://github.com/cloudflare/cf-terraforming/releases) or
build the Go source.

## Importing with Terraform state

`cf-terraforming` has the ability to generate the configuration for you to import
existing resources.

Depending on your version of Terraform, you can generate the `import` block
(Terraform 1.5+) using the `--modern-import-block` flag or the `terraform import`
compatible CLI output (all versions).

This command assumes you have already ran `cf-terraforming generate ...` to
output your resources.

```
# All versions of Terraform
cf-terraforming import \
  --resource-type "cloudflare_record" \
  --email $CLOUDFLARE_EMAIL \
  --key $CLOUDFLARE_API_KEY \
  --zone $CLOUDFLARE_ZONE_ID
```

```
# Terraform 1.5+ only
cf-terraforming import \
  --resource-type "cloudflare_record" \
  --modern-import-block \
  --email $CLOUDFLARE_EMAIL \
  --key $CLOUDFLARE_API_KEY \
  --zone $CLOUDFLARE_ZONE_ID
```

## Using non-standard registries

By default, we use the Hashicorp registry (registry.terraform.io) for looking up
the provider to introspect the schema. If you are attempting to use another
registry, you will need to provide the `--provider-registry-hostname` flag or
`CLOUDFLARE_PROVIDER_REGISTRY_HOSTNAME` environment variable to query the correct
registry.

## Using non-standard Terraform binaries

Internally, we use [`terraform-exec`](https://github.com/hashicorp/terraform-exec)
library to run Terraform operations in the same way that the CLI tooling would.
If a `terraform` binary is not available on your system path, we will attempt
to download the latest to use it.

Should you have the binary stored in a non-standard location, want to use an
existing binary, or you wish to provide a Terraform compatible binary (such as
`tofu`), you need to provide the `--terraform-binary-path` flag or
`CLOUDFLARE_TERRAFORM_BINARY_PATH` environment variable to instruct
`cf-terraforming` which you expect to use.

## Supported Resources

### v5

Any resource that is released within the Terraform Provider is automatically
supported for generation. Import support is not yet implemented.

### v4

Any resources not listed are currently not supported.

| Resource                                                                                                                                         | Resource Scope  | Generate Supported | Import Supported |
| ------------------------------------------------------------------------------------------------------------------------------------------------ | --------------- | ------------------ | ---------------- |
| [cloudflare_access_application](https://www.terraform.io/docs/providers/cloudflare/r/access_application)                                         | Account         | ✅                 | ✅               |
| [cloudflare_access_group](https://www.terraform.io/docs/providers/cloudflare/r/access_group)                                                     | Account         | ✅                 | ✅               |
| [cloudflare_access_identity_provider](https://www.terraform.io/docs/providers/cloudflare/r/access_identity_provider)                             | Account         | ✅                 | ❌               |
| [cloudflare_access_mutual_tls_certificate](https://www.terraform.io/docs/providers/cloudflare/r/access_mutual_tls_certificate)                   | Account         | ✅                 | ❌               |
| [cloudflare_access_policy](https://www.terraform.io/docs/providers/cloudflare/r/access_policy)                                                   | Account         | ❌                 | ❌               |
| [cloudflare_access_rule](https://www.terraform.io/docs/providers/cloudflare/r/access_rule)                                                       | Account         | ✅                 | ✅               |
| [cloudflare_access_service_token](https://www.terraform.io/docs/providers/cloudflare/r/access_service_token)                                     | Account         | ✅                 | ❌               |
| [cloudflare_account_member](https://www.terraform.io/docs/providers/cloudflare/r/account_member)                                                 | Account         | ✅                 | ✅               |
| [cloudflare_api_shield](https://www.terraform.io/docs/providers/cloudflare/r/api_shield)                                                         | Zone            | ✅                 | ❌               |
| [cloudflare_api_token](https://www.terraform.io/docs/providers/cloudflare/r/api_token)                                                           | User            | ❌                 | ❌               |
| [cloudflare_argo](https://www.terraform.io/docs/providers/cloudflare/r/argo)                                                                     | Zone            | ✅                 | ✅               |
| [cloudflare_authenticated_origin_pulls](https://www.terraform.io/docs/providers/cloudflare/r/authenticated_origin_pulls)                         | Zone            | ❌                 | ❌               |
| [cloudflare_authenticated_origin_pulls_certificate](https://www.terraform.io/docs/providers/cloudflare/r/authenticated_origin_pulls_certificate) | Zone            | ❌                 | ❌               |
| [cloudflare_bot_management](https://registry.terraform.io/providers/cloudflare/cloudflare/latest/docs/resources/bot_management)                  | Zone            | ✅                 | ✅               |
| [cloudflare_byo_ip_prefix](https://www.terraform.io/docs/providers/cloudflare/r/byo_ip_prefix)                                                   | Account         | ✅                 | ✅               |
| [cloudflare_certificate_pack](https://www.terraform.io/docs/providers/cloudflare/r/certificate_pack)                                             | Zone            | ✅                 | ✅               |
| [cloudflare_custom_hostname](https://www.terraform.io/docs/providers/cloudflare/r/custom_hostname)                                               | Zone            | ✅                 | ✅               |
| [cloudflare_custom_hostname_fallback_origin](https://www.terraform.io/docs/providers/cloudflare/r/custom_hostname_fallback_origin)               | Account         | ✅                 | ❌               |
| [cloudflare_custom_pages](https://www.terraform.io/docs/providers/cloudflare/r/custom_pages)                                                     | Account or Zone | ✅                 | ✅               |
| [cloudflare_custom_ssl](https://www.terraform.io/docs/providers/cloudflare/r/custom_ssl)                                                         | Zone            | ✅                 | ✅               |
| [cloudflare_filter](https://www.terraform.io/docs/providers/cloudflare/r/filter)                                                                 | Zone            | ✅                 | ✅               |
| [cloudflare_firewall_rule](https://www.terraform.io/docs/providers/cloudflare/r/firewall_rule)                                                   | Zone            | ✅                 | ✅               |
| [cloudflare_healthcheck](https://www.terraform.io/docs/providers/cloudflare/r/healthcheck)                                                       | Zone            | ✅                 | ✅               |
| [cloudflare_ip_list](https://www.terraform.io/docs/providers/cloudflare/r/ip_list)                                                               | Account         | ❌                 | ✅               |
| [cloudflare_list](https://www.terraform.io/docs/providers/cloudflare/r/list)                                                                     | Account         | ✅                 | ❌               |
| [cloudflare_load_balancer](https://www.terraform.io/docs/providers/cloudflare/r/load_balancer)                                                   | Zone            | ✅                 | ✅               |
| [cloudflare_load_balancer_monitor](https://www.terraform.io/docs/providers/cloudflare/r/load_balancer_monitor)                                   | Account         | ✅                 | ✅               |
| [cloudflare_load_balancer_pool](https://www.terraform.io/docs/providers/cloudflare/r/load_balancer_pool)                                         | Account         | ✅                 | ✅               |
| [cloudflare_logpull_retention](https://www.terraform.io/docs/providers/cloudflare/r/logpull_retention)                                           | Zone            | ❌                 | ❌               |
| [cloudflare_logpush_job](https://www.terraform.io/docs/providers/cloudflare/r/logpush_job)                                                       | Zone            | ✅                 | ❌               |
| [cloudflare_logpush_ownership_challenge](https://www.terraform.io/docs/providers/cloudflare/r/logpush_ownership_challenge)                       | Zone            | ❌                 | ❌               |
| [cloudflare_magic_firewall_ruleset](https://www.terraform.io/docs/providers/cloudflare/r/magic_firewall_ruleset)                                 | Account         | ❌                 | ❌               |
| [cloudflare_origin_ca_certificate](https://www.terraform.io/docs/providers/cloudflare/r/origin_ca_certificate)                                   | Zone            | ✅                 | ✅               |
| [cloudflare_page_rule](https://www.terraform.io/docs/providers/cloudflare/r/page_rule)                                                           | Zone            | ✅                 | ✅               |
| [cloudflare_rate_limit](https://www.terraform.io/docs/providers/cloudflare/r/rate_limit)                                                         | Zone            | ✅                 | ✅               |
| [cloudflare_record](https://www.terraform.io/docs/providers/cloudflare/r/record)                                                                 | Zone            | ✅                 | ✅               |
| [cloudflare_ruleset](https://www.terraform.io/docs/providers/cloudflare/r/ruleset)                                                               | Account or Zone | ✅                 | ✅               |
| [cloudflare_spectrum_application](https://www.terraform.io/docs/providers/cloudflare/r/spectrum_application)                                     | Zone            | ✅                 | ✅               |
| [cloudflare_tiered_cache](https://www.terraform.io/docs/providers/cloudflare/r/tiered_cache)                                                     | Zone            | ✅                 | ❌               |
| [cloudflare_teams_list](https://www.terraform.io/docs/providers/cloudflare/r/teams_list)                                                         | Account         | ✅                 | ✅               |
| [cloudflare_teams_location](https://www.terraform.io/docs/providers/cloudflare/r/teams_location)                                                 | Account         | ✅                 | ✅               |
| [cloudflare_teams_proxy_endpoint](https://www.terraform.io/docs/providers/cloudflare/r/teams_proxy_endpoint)                                     | Account         | ✅                 | ✅               |
| [cloudflare_teams_rule](https://www.terraform.io/docs/providers/cloudflare/r/teams_rule)                                                         | Account         | ✅                 | ✅               |
| [cloudflare_tunnel](https://www.terraform.io/docs/providers/cloudflare/r/tunnel)                                                                 | Account         | ✅                 | ✅               |
| [cloudflare_turnstile_widget](https://registry.terraform.io/providers/cloudflare/cloudflare/latest/docs/resources/turnstile_widget)              | Account         | ✅                 | ✅               |
| [cloudflare_url_normalization_settings](https://www.terraform.io/docs/providers/cloudflare/r/url_normalization_settings)                         | Zone            | ✅                 | ❌               |
| [cloudflare_waf_group](https://www.terraform.io/docs/providers/cloudflare/r/waf_group)                                                           | Zone            | ❌                 | ❌               |
| [cloudflare_waf_override](https://www.terraform.io/docs/providers/cloudflare/r/waf_override)                                                     | Zone            | ✅                 | ✅               |
| [cloudflare_waf_package](https://www.terraform.io/docs/providers/cloudflare/r/waf_package)                                                       | Zone            | ✅                 | ❌               |
| [cloudflare_waf_rule](https://www.terraform.io/docs/providers/cloudflare/r/waf_rule)                                                             | Zone            | ❌                 | ❌               |
| [cloudflare_waiting_room](https://www.terraform.io/docs/providers/cloudflare/r/waiting_room)                                                     | Zone            | ✅                 | ✅               |
| [cloudflare_worker_cron_trigger](https://www.terraform.io/docs/providers/cloudflare/r/worker_cron_trigger)                                       | Account         | ❌                 | ❌               |
| [cloudflare_worker_route](https://www.terraform.io/docs/providers/cloudflare/r/worker_route)                                                     | Zone            | ✅                 | ✅               |
| [cloudflare_worker_script](https://www.terraform.io/docs/providers/cloudflare/r/worker_script)                                                   | Account         | ❌                 | ❌               |
| [cloudflare_workers_kv](https://www.terraform.io/docs/providers/cloudflare/r/workers_kv)                                                         | Account         | ❌                 | ❌               |
| [cloudflare_workers_kv_namespace](https://www.terraform.io/docs/providers/cloudflare/r/workers_kv_namespace)                                     | Account         | ✅                 | ✅               |
| [cloudflare_zone](https://www.terraform.io/docs/providers/cloudflare/r/zone)                                                                     | Account         | ✅                 | ✅               |
| [cloudflare_zone_dnssec](https://www.terraform.io/docs/providers/cloudflare/r/zone_dnssec)                                                       | Zone            | ❌                 | ❌               |
| [cloudflare_zone_lockdown](https://www.terraform.io/docs/providers/cloudflare/r/zone_lockdown)                                                   | Zone            | ✅                 | ✅               |
| [cloudflare_zone_settings_override](https://www.terraform.io/docs/providers/cloudflare/r/zone_settings_override)                                 | Zone            | ✅                 | ❌               |

## Testing

To ensure changes don't introduce regressions this tool uses an automated test
suite consisting of HTTP mocks via go-vcr and Terraform configuration files to
assert against. The premise is that we mock the HTTP responses from the
Cloudflare API to ensure we don't need to create and delete real resources to
test. The Terraform files then allow us to build what the resource structure is
expected to look like and once the tool parses the API response, we can compare
that to the static file.

Suggested local testing steps:

1. Create a file with the basic provider configuration (do not commit this file)

```bash
cat > main.tf <<EOF
terraform {
  required_providers {
    cloudflare = {
      source = "cloudflare/cloudflare"
      version = "~> 4"
    }
  }
}
EOF
```

2. Initialize terraform

```bash
terraform init
```

3. Run tests (Cloudflare Install path should be path to repository)

```bash
make test
```

If you want to run a specific test case you can do so with the TESTARGS variable and -run flag

```bash
TESTARGS="-run '^TestResourceGeneration/cloudflare_teams_list'" make test
```

## Updating VCR cassettes

Periodically, it is a good idea to recreate the VCR cassettes used in our
testing to ensure they haven't drifted from actual responses. To do this, you
will need to:

- Create the appropriate resource in a Cloudflare account/zone you have access
  to. This is required as overwriting cassettes makes real API requests on your
  behalf.
- Invoke the test suite with `OVERWRITE_VCR_CASSETTES=true`,
  `CLOUDFLARE_DOMAIN=<real domain here>`, authentication credentials
  (`CLOUDFLARE_EMAIL`, `CLOUDFLARE_KEY`, `CLOUDFLARE_API_TOKEN`) and the test
  you want to update.
  Example of updating the DNS CAA record test with a zone I own:

```bash
  OVERWRITE_VCR_CASSETTES=true \
    CLOUDFLARE_DOMAIN="terraform.cfapi.net" \
    CLOUDFLARE_EMAIL="jb@example.com" \
    CLOUDFLARE_API_KEY="..." \
    TESTARGS="-run '^TestResourceGeneration/cloudflare_record_caa'"  \
    make test
```

- Commit your changes and push them via a Pull Request.
