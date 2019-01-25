package cmd

import (
	"fmt"
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go"

	"text/template"

	"github.com/spf13/cobra"
)

const spectrumApplicationTemplate = `
resource "cloudflare_spectrum_application" "{{.App.ID}}" {
    protocol = "{{.App.Protocol}}"
    dns = {
        type = "{{.App.DNS.Type}}"
        name = "{{.App.DNS.Name}}"
    }
{{if .App.OriginPort}}
    origin_port = "{{.App.OriginPort}}"
{{end}}
{{if .App.OriginDNS}}
    origin_dns = {
      name = "{{.App.OriginDNS.Name}}"
    }
{{end}}
{{if .App.IPFirewall}}
    ip_firewall = "{{.App.IPFirewall}}"
{{end}}
{{if .App.ProxyProtocol}}
    proxy_protocol = "{{.App.ProxyProtocol}}"
{{end}}
{{if .App.TLS}}
    tls = "{{.App.TLS}}"
{{end}}
    origin_direct = [{{range .App.OriginDirect}} "{{.}}", {{end}}]
}
`

func init() {
	rootCmd.AddCommand(spectrumApplicationCmd)
}

var spectrumApplicationCmd = &cobra.Command{
	Use:   "spectrum_application",
	Short: "Import a spectrum application into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		// Loop through all zones in account and fetch routes for each zone
		for _, zone := range zones {
			spectrumApplications, err := api.SpectrumApplications(zone.ID)

			if err != nil {
				// FIXME: api.SpectrumApplications can return 403 errors in the case
				// of permissions problems, for example, which will pollute tfstate
				fmt.Println(err)
				os.Exit(1)
			}

			if len(spectrumApplications) > 0 {
				for _, app := range spectrumApplications {
					spectrumAppParse(app)
				}
			}
		}
	},
}

func spectrumAppParse(app cloudflare.SpectrumApplication) {
	tmpl := template.Must(template.New("script").Funcs(templateFuncMap).Parse(spectrumApplicationTemplate))
	tmpl.Execute(os.Stdout,
		struct {
			App cloudflare.SpectrumApplication
		}{
			App: app,
		})
}
