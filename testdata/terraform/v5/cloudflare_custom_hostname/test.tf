resource "cloudflare_custom_hostname" "terraform_managed_resource" {
  hostname = "okwyujswsc.terraform.cfapi.net"
  zone_id  = "0da42c8d2132a9ddaf714f9e7c920711"
  ssl = {
    bundle_method         = "ubiquitous"
    certificate_authority = "google"
    id                    = "aa4bd600-4144-46a1-82b8-8ef525922877"
    method                = "txt"
    status                = "pending_validation"
    txt_name              = "_acme-challenge.okwyujswsc.terraform.cfapi.net"
    txt_value             = "28Rth8iqNkhKk8sPwpdkanLg5xMoVaEdiHsOzKHemUE"
    type                  = "dv"
  }
}
