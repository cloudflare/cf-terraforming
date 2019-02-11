# Cloudflare Terraforming
[![Go Report Card](https://goreportcard.com/badge/cloudflare/cf-terraforming)](https://goreportcard.com/report/cloudflare/cf-terraforming)

## Overview

cf-terraforming is a command line utility to facilitate terraforming your existing Cloudflare resources. It does this by using your account credentials to retrieve your configurations from the [Cloudflare API](https://api.cloudflare.com) and converting them to Terraform configurations that can be used with the [Terraform Cloudflare provider](https://www.terraform.io/docs/providers/cloudflare/index.html). 

This tool is ideal if you already have Cloudflare resources defined but want to start managing them via Terraform, and don't want to spend the time to manually write the Terraform configuration to describe them.

## Example

Running: 

```go run cmd/cf-terraforming/main.go --email <your-cloudflare-account-email> --key <your-cloudflare-api-key> -a <your-cloudflare-account-id> spectrum_application```

will contact the Cloudflare API on your behalf and result in a valid Terraform configuration representing the resource you requested:

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

## Prerequisites 
* A Cloudflare account with resources defined (e.g. a few zones, some load balancers, spectrum applications, etc)
* A valid Cloudflare API key and sufficient permissions to access the resources you are requesting via the API
* A working [installation of Go](https://golang.org/doc/install)

## Installation

```bash
$ go get -u github.com/cloudflare/cf-terraforming/...
```
This will fetch the cf-terraforming tool as well as its dependencies, updating them as necessary.

## Usage 

You can use ```go run``` to build and execute the binary in a single command like so: 

```
go run cmd/cf-terraforming/main.go --email <your-cloudflare-account-email> --key <your-cloudflare-api-key> -z <your-cloudflare-zone> -a <your-cloudflare-account-id> <resource>
```
where ```resource``` is one of the **supported resources**.

## Supported resources

* [access_application](https://www.terraform.io/docs/providers/cloudflare/r/access_application.html)
* [access_rule](https://www.terraform.io/docs/providers/cloudflare/r/access_rule.html)
* [account_member](https://www.terraform.io/docs/providers/cloudflare/r/account_member.html)
* [custom_pages](https://www.terraform.io/docs/providers/cloudflare/r/custom_pages.html)
* [filter](https://www.terraform.io/docs/providers/cloudflare/r/filter.html)
* [firewall_rule](https://www.terraform.io/docs/providers/cloudflare/r/firewall_rule.html)
* [load_balancer](https://www.terraform.io/docs/providers/cloudflare/r/load_balancer.html)
* [load_balancer_pool](https://www.terraform.io/docs/providers/cloudflare/r/load_balancer_pool.html)
* [rate_limit](https://www.terraform.io/docs/providers/cloudflare/r/rate_limit.html)
* [record](https://www.terraform.io/docs/providers/cloudflare/r/record.html)
* [spectrum_application](https://www.terraform.io/docs/providers/cloudflare/r/spectrum_application.html) 
* [waf_rule](https://www.terraform.io/docs/providers/cloudflare/r/waf_rule.html)
* [worker_route](https://www.terraform.io/docs/providers/cloudflare/r/worker_route.html)
* [worker_script](https://www.terraform.io/docs/providers/cloudflare/r/worker_script.html)
* [zone](https://www.terraform.io/docs/providers/cloudflare/r/zone.html) 
* [zone_lockdown](https://www.terraform.io/docs/providers/cloudflare/r/zone_lockdown.html)
* [zone_settings_override](https://www.terraform.io/docs/providers/cloudflare/r/zone_settings_override.html) 
