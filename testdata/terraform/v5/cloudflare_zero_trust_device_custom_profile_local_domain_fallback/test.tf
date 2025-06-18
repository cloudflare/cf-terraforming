# resource "cloudflare_zero_trust_device_custom_profile_local_domain_fallback" "terraform_managed_resource" {
#   account_id = "699d98642c564d2e855e9661899b7252"
#   policy_id = "f174e90a-fafe-4643-bbbc-4a0ed4fc8415"
#   domains = [{
#     suffix = "example.com"
#     description = "Domain bypass for local development"
#     dns_server = ["1.1.1.1"]
#   }]
# }
resource "cloudflare_zero_trust_device_custom_profile_local_domain_fallback" "test" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  policy_id  = "61816005-1544-4f75-8475-d6258cdede81"
  domains = [
    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "16.172.in-addr.arpa"
    },
    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "168.192.in-addr.arpa"
    },
    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "17.172.in-addr.arpa"
    },
    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "18.172.in-addr.arpa"
    },
    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "10.in-addr.arpa"
    },
    {
      suffix = "home.arpa"
    },
    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "maninvestments.com"
    },
    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "man.com"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "19.172.in-addr.arpa"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "20.172.in-addr.arpa"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "21.172.in-addr.arpa"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "22.172.in-addr.arpa"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "23.172.in-addr.arpa"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "24.172.in-addr.arpa"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "25.172.in-addr.arpa"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "26.172.in-addr.arpa"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "27.172.in-addr.arpa"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "28.172.in-addr.arpa"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "29.172.in-addr.arpa"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "30.172.in-addr.arpa"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "31.172.in-addr.arpa"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "adm.rmcm2.reuters.net"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "ahl"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "ahl.com"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "collab.reuasmb.net"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "collab.reuters.net"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "corp"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "corp.silverminecap.com"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "cp.reutest.net"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "database.windows.net"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "extranet.reured.biz"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "extranet.reuters.biz"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "extranet.reutest.biz"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "extranet.thomsonreuters.biz"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "friskman.com"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "frm"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "frmhedge.com"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "fxall.com"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "fxtrading.reublue.com"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "glg"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "glgpartners.com"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "gpm"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "m"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "man.com"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "maninvestments.ad.man.com"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "maninvestments.com"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "ms.crd.com"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "num"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "numeric.com"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "privatelink.blob.core.windows.net"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "qarl"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "ra.lcl"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "reuters.com"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "rservices.com"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "rtextrading.reublue.com"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "servicebus.windows.net"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "sip.reuters.net"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "thomsonreuters.com"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "thomsonreuters.net"
    },

    {
      dns_server = [
        "10.192.207.245",
        "10.196.207.245",
      ]
      suffix = "westeurope.azmk8s.io"
    },
  ]
}
