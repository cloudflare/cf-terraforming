package cmd

import (
	"os"
	"strings"

	cloudflare "github.com/cloudflare/cloudflare-go"

	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const spectrumApplicationTemplate = `
resource "cloudflare_spectrum_application" "spectrum_application_{{.App.ID}}" {
    zone_id = "{{.Zone.ID}}"
    protocol = "{{.App.Protocol}}"
    dns {
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
		log.Debug("Importing Spectrum application data")
		// Loop through all zones in account and fetch routes for each zone
		for _, zone := range zones {
			spectrumApplications, err := api.SpectrumApplications(zone.ID)

			if err != nil {
				if strings.Contains(err.Error(), "HTTP status 403") {
					log.WithFields(logrus.Fields{
						"ID": zone.ID,
					}).Debug("Insufficient permissions for accessing zone")
					continue
				}
				log.Debug(err)
			}

			if len(spectrumApplications) > 0 {
				for _, app := range spectrumApplications {

					log.WithFields(logrus.Fields{
						"ID": app.ID,
					}).Debug("Processing spectrum app")

					if tfstate {
						// TODO: Implement state dump
					} else {
						spectrumAppParse(app, zone)
					}
				}
			}
		}
	},
}

func spectrumAppParse(app cloudflare.SpectrumApplication, zone cloudflare.Zone) {
	// modified this section to support zone id
	tmpl := template.Must(template.New("script").Funcs(templateFuncMap).Parse(spectrumApplicationTemplate))
	tmpl.Execute(os.Stdout,
		struct {
			Zone cloudflare.Zone
			App cloudflare.SpectrumApplication
		}{
			App: app,
			Zone: zone,
		})
}
