package cmd

import (
	"fmt"
	"os"

	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const zoneLockdownTemplate = `
resource "cloudflare_zone_lockdown" "{{replace .Zone.Name "." "_"}}_{{.Lockdown.ID}}" {
    zone_id = "{{.Zone.ID}}"
    description = "{{.Lockdown.Description}}"
    urls = [
{{range .Lockdown.URLs}}
        "{{.}}",
{{end}}
    ]
    configurations = [
{{range .Lockdown.Configurations}}
        {
            target = "{{.Target}}"
            value = "{{.Value}}"
        },
{{end}}
    ]
}
`

func init() {
	rootCmd.AddCommand(zoneLockdownCmd)
}

var zoneLockdownCmd = &cobra.Command{
	Use:   "zone_lockdown",
	Short: "Import Zone Lockdown data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Importing Zone Lockdown data")

		for _, zone := range zones {
			log.WithFields(logrus.Fields{
				"ID":   zone.ID,
				"Name": zone.Name,
			}).Debug("Processing zone")

			totalPages := 999

			for page := 1; page <= totalPages; page++ {
				lockdowns, err := api.ListZoneLockdowns(zone.ID, page)

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				totalPages = lockdowns.TotalPages

				for _, r := range lockdowns.Result {

					log.WithFields(logrus.Fields{
						"ID":  r.ID,
						"URL": r.URLs,
					}).Debug("Processing lockdown")

					zoneLockdownParse(zone, r)
				}
			}
		}
	},
}

func zoneLockdownParse(zone cloudflare.Zone, lockdown cloudflare.ZoneLockdown) {
	tmpl := template.Must(template.New("zone_lockdown").Funcs(templateFuncMap).Parse(zoneLockdownTemplate))
	if err := tmpl.Execute(os.Stdout,
		struct {
			Zone     cloudflare.Zone
			Lockdown cloudflare.ZoneLockdown
		}{
			Zone:     zone,
			Lockdown: lockdown,
		}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
