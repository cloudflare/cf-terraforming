resource "cloudflare_zero_trust_tunnel_cloudflared_config" "terraform_managed_resource" {
  account_id = "f037e56e89293a057740de681ac9abbe"
  source     = "cloudflare"
  tunnel_id  = "285f508d-d6ef-4ce4-9293-983d5bdc269e"
  config = {
    ingress = [{
      hostname = "foo"
      originRequest = {
        bastionMode            = false
        caPool                 = ""
        connectTimeout         = 15
        disableChunkedEncoding = false
        http2Origin            = false
        httpHostHeader         = ""
        keepAliveConnections   = 100
        keepAliveTimeout       = 90
        noHappyEyeballs        = false
        noTLSVerify            = false
        originServerName       = ""
        proxyAddress           = "127.0.0.1"
        proxyPort              = 0
        proxyType              = ""
        tcpKeepAlive           = 30
        tlsTimeout             = 10
      }
      path    = "/bar"
      service = "http://10.0.0.2:8080"
      }, {
      hostname = "bar"
      originRequest = {
        access = {
          audTag   = ["AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"]
          required = true
          teamName = "terraform"
        }
        bastionMode            = false
        caPool                 = ""
        connectTimeout         = 30
        disableChunkedEncoding = false
        http2Origin            = false
        httpHostHeader         = ""
        keepAliveConnections   = 100
        keepAliveTimeout       = 90
        noHappyEyeballs        = false
        noTLSVerify            = false
        originServerName       = ""
        proxyAddress           = "127.0.0.1"
        proxyPort              = 0
        proxyType              = ""
        tcpKeepAlive           = 30
        tlsTimeout             = 10
      }
      path    = "/foo"
      service = "http://10.0.0.3:8081"
      }, {
      service = "https://10.0.0.4:8082"
    }]
    originRequest = {
      bastionMode            = false
      caPool                 = "/path/to/unsigned/ca/pool"
      connectTimeout         = 60
      disableChunkedEncoding = false
      http2Origin            = true
      httpHostHeader         = "baz"
      ipRules = [{
        ports  = [80, 443]
        prefix = "/web"
      }]
      keepAliveConnections = 1024
      keepAliveTimeout     = 60
      noHappyEyeballs      = false
      noTLSVerify          = false
      originServerName     = "foobar"
      proxyAddress         = "10.0.0.1"
      proxyPort            = 8123
      proxyType            = "socks"
      tcpKeepAlive         = 60
      tlsTimeout           = 60
    }
    warp-routing = {
      enabled = true
    }
  }
}

