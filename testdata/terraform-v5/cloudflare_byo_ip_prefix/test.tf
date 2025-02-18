resource "cloudflare_byo_ip_prefix" "example_byo_ip_prefix" {
  account_id = "258def64c72dae45f3e4c8516e2111f2"
  asn = 209242
  cidr = "192.0.2.0/24"
  loa_document_id = "d933b1530bc56c9953cf8ce166da8004"
}
