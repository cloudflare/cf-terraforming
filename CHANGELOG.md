## Unreleased 

- deps: update hashicorp/terraform-json to 0.11.0 (#261)
- generate: remove `cloudflare_zone.status` from output
- generate: remap `cloudflare_load_balancer` `default_pool_ids` and `fallback_pool_id` to their schema values

## 0.3.0 (2021-05-20)

- docs: update Terraform registry documentation links (#255)
- generate: add support for page rules (#259)

## 0.2.0 (2021-04-28)

- generate: handle non-string values for IDs ([#249](https://github.com/cloudflare/cf-terraforming/issues/249))
- deps: bump zclconf/go-cty to 1.8.2 ([#247](https://github.com/cloudflare/cf-terraforming/issues/247))
- deps: update hashicorp/terraform-exec to 0.13.3 to address GPG revocation by HashiCorp ([HCSEC-2021-12](https://discuss.hashicorp.com/t/hcsec-2021-12-codecov-security-event-and-hashicorp-gpg-key-exposure/23512)) ([#250](https://github.com/cloudflare/cf-terraforming/issues/250))
- deps: add explicit dependency for hashicorp/go-getter ([#253](https://github.com/cloudflare/cf-terraforming/issues/253))

## 0.1.1 (2021-04-15)

- generate: remove `tfexec.LockTimeout` on init for Terraform 0.15 support

## 0.1.0 (2021-04-13)

- Revamped internals to support dynamic generation of resources.
