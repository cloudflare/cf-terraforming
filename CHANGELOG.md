## Unreleased

- generate: handle non-string values for IDs (#249)
- deps: bump zclconf/go-cty to 1.8.2 (#247)
- deps: update hashicorp/terraform-exec to 0.13.3 to address GPG revocation by HashiCorp ([HCSEC-2021-12](https://discuss.hashicorp.com/t/hcsec-2021-12-codecov-security-event-and-hashicorp-gpg-key-exposure/23512)) (#250)
- deps: add explicit dependency for hashicorp/go-getter (#253)

## 0.1.1 (2021-04-15)

- generate: remove `tfexec.LockTimeout` on init for Terraform 0.15 support

## 0.1.0 (2021-04-13)

- Revamped internals to support dynamic generation of resources.
