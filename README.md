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

> NOTE: If you would like to export resources compatible with Terraform < 0.12.x,
> you will need to download an older release as this tool no longer supports it.

## Usage

```
Usage:
  cf-terraforming [command]

Available Commands:
  generate    Fetch resources from the Cloudflare API and generate the respective Terraform stanzas
  help        Help about any command
  import      Output `terraform import` compatible commands in order to import resources into state
  version     Print the version number of cf-terraforming

Flags:
  -a, --account string         Use specific account ID for commands
  -c, --config string          Path to configuration file (default is $HOME/.cf-terraforming.yaml)
  -e, --email string           API Email address associated with your account
  -h, --help                   Help for cf-terraforming
  -k, --key string             API Key generated on the 'My Profile' page. See: https://dash.cloudflare.com/profile
      --resource-type string   Which resource you wish to generate
  -t, --token string           API Token
  -v, --verbose                Specify verbose output (same as setting log level to debug)
  -z, --zone string            Limit the export to a single zone ID

Use "cf-terraforming [command] --help" for more information about a command.
```

## Authentication

Cloudflare supports two authentication methods to the API:
* API Token - gives access only to resources and permissions specified for that token (recommended)
* API key - gives access to everything your user profile has access to

Both can be retrieved on [profile page](https://dash.cloudflare.com/profile/api-tokens).

**A note on storing your credentials securely:** We recommend that you store
your Cloudflare credentials (API key, email, token) as environment variables as
demonstrated below.

```bash
# if using API Token
export CLOUDFLARE_API_TOKEN='Hzsq3Vub-7Y-hSTlAaLH3Jq_YfTUOCcgf22_Fs-j'

# if using API Key
export CLOUDFLARE_EMAIL='user@example.com'
export CLOUDFLARE_API_KEY='1150bed3f45247b99f7db9696fffa17cbx9'

# specify Account ID
export CLOUDFLARE_ACCOUNT_ID='81b06ss3228f488fh84e5e993c2dc17'

# now call cf-terraforming, e.g.
cf-terraforming generate --resource-type "cloudflare_record" --account $CLOUDFLARE_ACCOUNT_ID
```

cf-terraforming supports the following environment variables:
* CLOUDFLARE_API_TOKEN - API Token based authentication
* CLOUDFLARE_EMAIL, CLOUDFLARE_API_KEY - API Key based authentication

## Example usage

```bash
$ cf-terraforming generate --account $CLOUDFLARE_ACCOUNT_ID --resource-type "cloudflare_record"
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

* A Cloudflare account with resources defined (e.g. a few zones, some load
  balancers, spectrum applications, etc)
* A valid Cloudflare API key and sufficient permissions to access the resources
  you are requesting via the API
* A working [installation of Go](https://golang.org/doc/install) at least
  v1.15.x.

## Installation

If you use Homebrew on MacOS, you can run the following:

```bash
brew tap cloudflare/cloudflare
brew install cloudflare/cloudflare/cf-terrarforming
```

Otherwise:

```bash
$ GO111MODULE=on go get -u github.com/cloudflare/cf-terraforming/...
```

This will fetch the `cf-terraforming` tool as well as its dependencies, updating
them as necessary, build and install the package in your `$GOPATH` (usually
`~/go/bin`). You can check your current GOPATH by running:

```bash
$ go env | grep GOPATH
```

## Importing with Terraform state

As of the latest release, `cf-terraforming` will output the `terraform import`
compatible commands for you when you invoke the `import` command. This command
assumes you have already ran `cf-terraforming generate ...` to output your
resources.

In the future we aim to automate this however for now, it is a manual step to
allow flexibility in directory structure.

```
$ cf-terraforming import --resource-type "cloudflare_record" --email $CLOUDFLARE_EMAIL --key $CLOUDFLARE_API_KEY -z "example.com"
```

## Testing

To ensure changes don't introduce regressions this tool uses an automated test
suite consisting of HTTP mocks via go-vcr and Terraform configuration files to
assert against. The premise is that we mock the HTTP responses from the
Cloudflare API to ensure we don't need to create and delete real resources to
test. The Terraform files then allow us to build what the resource structure is
expected to look like and once the tool parses the API response, we can compare
that to the static file.
