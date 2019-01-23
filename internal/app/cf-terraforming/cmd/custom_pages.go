package cmd

import (
	"fmt"
	"log"
	"os"
	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/spf13/cobra"
)

const customPagesTemplate = `
resource "cloudflare_custom_pages" "{{.CustomPage.ID}}" {
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
		log.Print("Importing Custom Pages data")

		for _, zone := range zones {
			log.Printf("[DEBUG] Processing zone: ID %s, Name %s", zone.ID, zone.Name)

			customPages, err := api.CustomPages(&cloudflare.CustomPageOptions{ZoneID: zone.ID})

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for _, r := range customPages {
				log.Printf("[DEBUG] Custom Page ID %s, URL %s, Description %s\n", r.ID, r.URL, r.Description)
				customPagesParse(zone, r)
			}
		}
	},
}

func customPagesParse(zone cloudflare.Zone, customPage cloudflare.CustomPage) {
	tmpl := template.Must(template.New("custom_pages").Funcs(templateFuncMap).Parse(customPagesTemplate))
	tmpl.Execute(os.Stdout,
		struct {
			Zone       cloudflare.Zone
			CustomPage cloudflare.CustomPage
		}{
			Zone:       zone,
			CustomPage: customPage,
		})
}
