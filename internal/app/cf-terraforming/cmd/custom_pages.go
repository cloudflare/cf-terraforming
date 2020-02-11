package cmd

import (
	"os"
	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const customPagesTemplate = `
resource "cloudflare_custom_pages" "custom_pages_{{.CustomPage.ID}}" {
  zone_id = "{{.Zone.ID}}"
  type    = "{{.CustomPage.ID}}"
  {{if .CustomPage.URL}}url     = "{{.CustomPage.URL}}"{{end}}
  state   = "{{.CustomPage.State}}"
}
`

func init() {
	rootCmd.AddCommand(customPagesCmd)
}

var customPagesCmd = &cobra.Command{
	Use:   "custom_pages",
	Short: "Import Custom Pages data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Importing Custom Pages data")

		for _, zone := range zones {

			log.WithFields(logrus.Fields{
				"ID":   zone.ID,
				"Name": zone.Name,
			}).Debug("Processing zone")

			customPages, err := api.CustomPages(&cloudflare.CustomPageOptions{ZoneID: zone.ID})

			if err != nil {
				log.Error(err)
				return
			}

			for _, r := range customPages {

				log.WithFields(logrus.Fields{
					"ID":          r.ID,
					"URL":         r.URL,
					"Description": r.Description,
				}).Debug("Processing custom page")

				if tfstate {
					// TODO: Implement state dump
				} else {
					customPagesParse(zone, r)
				}
			}
		}
	},
}

func customPagesParse(zone cloudflare.Zone, customPage cloudflare.CustomPage) {
	tmpl := template.Must(template.New("custom_pages").Funcs(templateFuncMap).Parse(customPagesTemplate))
	err := tmpl.Execute(os.Stdout,
		struct {
			Zone       cloudflare.Zone
			CustomPage cloudflare.CustomPage
		}{
			Zone:       zone,
			CustomPage: customPage,
		})
	if err != nil {
		log.Error(err)
	}
}
