## Unreleased

## 0.6.0 (2021-12-14)

- generate: add support for nested map interfaces in request ([#332](https://github.com/cloudflare/cf-terraforming/pull/332))
- generate: add support for CLOUDFLARE_API_HOSTNAME ([#332](https://github.com/cloudflare/cf-terraforming/pull/332))

## 0.5.0 (2021-10-26)

- generate: add support for `cloudflare_ruleset` ([#305](https://github.com/cloudflare/cf-terraforming/issues/305))
- generate: don't export "universal" `cloudflare_certificate_pack` ([#306](https://github.com/cloudflare/cf-terraforming/issues/306))
- generate: don't export "mobile_redirect" `cloudflare_zone_settings_override` ([#320](https://github.com/cloudflare/cf-terraforming/issues/320))

## 0.4.0 (2021-09-01)

- generate: better validation that `--resource-type` is provided ([#302](https://github.com/cloudflare/cf-terraforming/issues/302))
- import: fix the `:id` parameter of `cloudflare_zone` ([#300](https://github.com/cloudflare/cf-terraforming/issues/300))
- generate: add the ability to export PTR DNS records for `cloudflare_record` ([#299](https://github.com/cloudflare/cf-terraforming/issues/299))
- generate: add support for `cloudflare_zone_settings_override` ([#298](https://github.com/cloudflare/cf-terraforming/issues/298))
- generate: add support for `cloudflare_rate_limit` ([#296](https://github.com/cloudflare/cf-terraforming/issues/296))
- generate: add support for `cloudflare_argo` ([#295](https://github.com/cloudflare/cf-terraforming/issues/295))
- deps: update github.com/spf13/cobra from 1.1.3 to 1.2.1 ([#277](https://github.com/cloudflare/cf-terraforming/issues/277))
- deps: update hashicorp/terraform-exec from 0.13.3 to 0.14.0 ([#274](https://github.com/cloudflare/cf-terraforming/issues/274))
- deps: update cloudflare/cloudflare-go from 0.17.0 to 0.21.0 ([#271](https://github.com/cloudflare/cf-terraforming/issues/271), [#283](https://github.com/cloudflare/cf-terraforming/issues/283), [#287](https://github.com/cloudflare/cf-terraforming/issues/287), [#293](https://github.com/cloudflare/cf-terraforming/issues/293))
- deps: update zclconf/go-cty from 1.8.3 to 1.8.4  ([#270](https://github.com/cloudflare/cf-terraforming/issues/270))
- deps: update spf13/viper from 1.7.1 to 1.9.1 ([#267](https://github.com/cloudflare/cf-terraforming/issues/249), [#273](https://github.com/cloudflare/cf-terraforming/issues/273), [#278](https://github.com/cloudflare/cf-terraforming/issues/278), [#290](https://github.com/cloudflare/cf-terraforming/issues/290))
- deps: update hashicorp/terraform-json to 0.12.0 ([#261](https://github.com/cloudflare/cf-terraforming/issues/261), [#272](https://github.com/cloudflare/cf-terraforming/issues/272))
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
