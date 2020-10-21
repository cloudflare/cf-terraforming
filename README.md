# Cloudflare Terraforming
[![Go Report Card](https://goreportcard.com/badge/cloudflare/cf-terraforming)](https://goreportcard.com/report/cloudflare/cf-terraforming)

## Overview

cf-terraforming is a command line utility to facilitate terraforming your existing Cloudflare resources. It does this by using your account credentials to retrieve your configurations from the [Cloudflare API](https://api.cloudflare.com) and converting them to Terraform configurations that can be used with the [Terraform Cloudflare provider](https://www.terraform.io/docs/providers/cloudflare/index.html).

This tool is ideal if you already have Cloudflare resources defined but want to start managing them via Terraform, and don't want to spend the time to manually write the Terraform configuration to describe them.

## Usage

```
Usage:
  cf-terraforming [command]

Available Commands:
  access_application     Import Access Application data into Terraform
  access_policy          Import Access Policy data into Terraform
  access_rule            Import Access Rule data into Terraform
  account_member         Import Account Member data into Terraform
  all                    Import all Cloudflare resources into Terraform
  custom_pages           Import Custom Pages data into Terraform
  filter                 Import Filter data into Terraform
  firewall_rule          Import Firewall Rule data into Terraform
  help                   Help about any command
  load_balancer          Import a load balancer into Terraform
  load_balancer_monitor  Import a load balancer monitor into Terraform
  load_balancer_pool     Import a load balancer pool into Terraform
  page_rule              Import Page Rule data into Terraform
  rate_limit             Import Rate Limit data into Terraform
  record                 Import Record data into Terraform
  spectrum_application   Import a spectrum application into Terraform
  version                Print the version number of cf-terraforming
  waf_rule               Import WAF Rule data into Terraform
  worker_route           Import a worker route into Terraform
  worker_script          Import a worker script into Terraform
  zone                   Import zone data into Terraform
  zone_lockdown          Import Zone Lockdown data into Terraform
  zone_settings_override Import Zone Settings Override data into Terraform

Flags:
  -a, --account string        Use specific account ID for import
  -c, --config string         config file (default is $HOME/.cf-terraforming.yaml)
  -e, --email string          API Email address associated with your account
  -h, --help                  help for cf-terraforming
  -k, --key string            API Key generated on the 'API Tokens' page, under 'API Keys'. See: https://dash.cloudflare.com/profile/api-tokens
  -l, --loglevel string       Specify logging level: (trace, debug, info, warn, error, fatal, panic)
  -o, --organization string   Use specific organization ID for import (deprecated, use -a instead)
  -s, --tfstate               Export tfstate for the given resource instead of HCL Terraform config (default ! See caveats below !)
  -t, --token string          API Token generated on the 'API Tokens' page, under 'API Tokens'. See: https://dash.cloudflare.com/profile/api-tokens
  -v, --verbose               Specify verbose output (same as setting log level to debug)
  -z, --zone string           Limit the export to a single zone (name or ID)

Use "cf-terraforming [command] --help" for more information about a command.
```

## Example

**A note on storing your credentials securely:** We recommend that you store your Cloudflare credentials (API key, email, account ID, etc) as environment variables as demonstrated below.

```cf-terraforming --email $CLOUDFLARE_EMAIL --key $CLOUDFLARE_TOKEN --account $CLOUDFLARE_ACCOUNT_ID spectrum_application```

will contact the Cloudflare API on your behalf and result in a valid Terraform configuration representing the **resource** you requested:

```
resource "cloudflare_spectrum_application" "1150bed3f45247b99f7db9696fffa17cbx9" {
    protocol = "tcp/8000"
    dns = {
        type = "CNAME"
        name = "example.com"
    }
    ip_firewall = "true"
    tls = "off"
    origin_direct = [ "tcp://37.241.37.138:8000", ]
}
```

See the currently **supported resources** below.

## Download all Cloudflare resources

Use the **all** command to download everything and convert it into Terraform config.

```cf-terraforming --email $CLOUDFLARE_EMAIL --key $CLOUDFLARE_TOKEN --account $CLOUDFLARE_ACCOUNT_ID all```

## Controlling output and verbose mode
By default, cf-terraforming will not output any log type messages to stdout when run, so as to not pollute your generated Terraform config files and to allow you to cleanly redirect cf-terraforming output to existing Terraform configs.

However, it can be useful when debugging issues to specify a logging level, like so:

```
cf-terraforming --email $CLOUDFLARE_EMAIL --key $CLOUDFLARE_TOKEN -a 1233455678d876bc764b5f763af7644411 -l="debug" spectrum_application

DEBU[0000] Initializing cloudflare-go                    API email=apicloudflare.com Account ID=e9e138b6x52ea331b359a2ddfc6a8 Organization ID= Zone name=example.com
DEBU[0000] Selecting zones for import
DEBU[0000] Zones selected:
DEBU[0000] Zone                                          ID=81b06ss3228f488fh84e5e993c2dc17 Name=example.com
DEBU[0000] Importing zone settings data
```

For convenience, you can set the verbose flag, which is functionally equivalent to setting a log level of debug:

```
cf-terraforming --email $CLOUDFLARE_EMAIL --key $CLOUDFLARE_TOKEN -a 1233455678d876bc764b5f763af7644411 -v spectrum_application
```

## Prerequisites
* A Cloudflare account with resources defined (e.g. a few zones, some load balancers, spectrum applications, etc)
* A valid Cloudflare API key and sufficient permissions to access the resources you are requesting via the API
* A working [installation of Go](https://golang.org/doc/install) at least v1.12.x.

## Installation

```bash
$ go get -u github.com/cloudflare/cf-terraforming/...
```
This will fetch the cf-terraforming tool as well as its dependencies, updating them as necessary, build and install the package in your `$GOPATH` (usually `~/go/bin`). You can check your current GOPATH by running:

```bash
$ go env | grep GOPATH
```

## Experimental: exporting tfstate

Work is underway to support downloading your resources via cf-terraform as valid tfstate, which will allow you to skip importing and let Terraform become aware of and begin managing your resources immediately.

Currently, only the worker_route command supports the --tfstate flag, but support for downloading tfstate will steadily be added across the other resources supported by cf-terraforming.

To use this currently experimental feature, pass the --tfstate (-s) flag to your command like so:

```
$ cf-terraforming --email $CLOUDFLARE_EMAIL --key $CLOUDFLARE_API_KEY -z example.com --account $CLOUDFLARE_ACCOUNT_ID --tfstate worker_route

```

The output will look something like this:

```
{
    "version": 1,
    "terraform_version": "",
    "serial": 0,
    "lineage": "",
    "modules": [
    {
        "path": ["root"],
        "depends_on": [],
        "outputs":
        {},
        "resources":
        {
            "cloudflare_worker_route.3489743985643594723e238ce":
            {
                "primary":
                {
                    "id": "249653298b6041148b411b4723e238ce",
                    "attributes":
                    {
                        "enabled": "true",
                        "id": "249653298b6041148b411b4723e238ce",
                        "multi_script": "true",
                        "pattern": "example.com",
                        "zone": "example.com",
                        "zone_id": "z1b06143shshs3223e5ec83c2z1klop"
                    },
                    "meta":
                    {},
                    "tainted": false
                },
                "depends_on": [],
                "deposed": [],
                "provider": "provider.cloudflare",
                "type": "cloudflare_worker_route"
            },
            "cloudflare_worker_route.zZ3exdfyt3e76274ef25dda436":
            {
                "primary":
                {
                    "id": "9b3e8f2cc874zsde9274ef25dda58f36",
                    "attributes":
                    {
                        "enabled": "true",
                        "id": "9b3e8f2cc87444fb9274ef25dda58f36",
                        "multi_script": "true",
                        "pattern": "example.com/*",
                        "zone": "example.com",
                        "zone_id": "z1b06143shshs3223e5ec83c2z1klop"
                    },
                    "meta":
                    {},
                    "tainted": false
                },
                "depends_on": [],
                "deposed": [],
                "provider": "provider.cloudflare",
                "type": "cloudflare_worker_route"
            }
        }
    }]
}
```

This means that when you next run ```$ terraform plan```, Terraform should see no difference between your defined resources and your current tfstate.

## Supported resources

The following resources can be downloaded into [Terraform HCL format](https://www.terraform.io/docs/configuration/syntax.html) right now. Some resources also support the `--tfstate` flag. Support across the remaining commands will be added over time.

| Resource | Generating HCL | Generating tfstate |
|----------|----------------|--------------------|
| [access_application](https://www.terraform.io/docs/providers/cloudflare/r/access_application.html) | ✔️ | ✔️ |
| [access_policy](https://www.terraform.io/docs/providers/cloudflare/r/access_policy.html) | ✔️ | ✔️ |
| [access_rule](https://www.terraform.io/docs/providers/cloudflare/r/access_rule.html) | ✔️ | ✔️ |
| [account_member](https://www.terraform.io/docs/providers/cloudflare/r/account_member.html) | ✔️ | ✖️ |
| [custom_pages](https://www.terraform.io/docs/providers/cloudflare/r/custom_pages.html) | ✔️ | ✖️ |
| [filter](https://www.terraform.io/docs/providers/cloudflare/r/filter.html) | ✔️ | ✔️ |
| [firewall_rule](https://www.terraform.io/docs/providers/cloudflare/r/firewall_rule.html) | ✔️ | ✔️ |
| [load_balancer](https://www.terraform.io/docs/providers/cloudflare/r/load_balancer.html) | ✔️ | ✖️ |
| [load_balancer_pool](https://www.terraform.io/docs/providers/cloudflare/r/load_balancer_pool.html) | ✔️ | ✖️ |
| [load_balancer_monitor](https://www.terraform.io/docs/providers/cloudflare/r/load_balancer_monitor.html) | ✔️ | ✖️ |
| [page_rule](https://www.terraform.io/docs/providers/cloudflare/r/page_rule.html) | ✔️ | ✖️ |
| [rate_limit](https://www.terraform.io/docs/providers/cloudflare/r/rate_limit.html) | ✔️ | ✔️ |
| [record](https://www.terraform.io/docs/providers/cloudflare/r/record.html) | ✔️ | ✔️ |
| [spectrum_application](https://www.terraform.io/docs/providers/cloudflare/r/spectrum_application.html) | ✔️ | ✖️ |
| [waf_rule](https://www.terraform.io/docs/providers/cloudflare/r/waf_rule.html) | ✔️ | ✖️ |
| [worker_route](https://www.terraform.io/docs/providers/cloudflare/r/worker_route.html) | ✔️ | ✔️ |
| [worker_script](https://www.terraform.io/docs/providers/cloudflare/r/worker_script.html) | ✔️ | ✖️ |
| [zone](https://www.terraform.io/docs/providers/cloudflare/r/zone.html) | ✔️ | ✔️ |
| [zone_lockdown](https://www.terraform.io/docs/providers/cloudflare/r/zone_lockdown.html) | ✔️ | ✔️ |
| [zone_settings_override](https://www.terraform.io/docs/providers/cloudflare/r/zone_settings_override.html) | ✔️ | ✖️ |
