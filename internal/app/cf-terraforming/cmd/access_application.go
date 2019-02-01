package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/spf13/cobra"
)

const accessApplicationTemplate = `
resource "cloudflare_access_application" "{{.App.ID}}" {
 	zone_id = "{{.Zone.ID}}"
 	name = "{{.App.Name}}"
 	domain = "{{.App.Domain}}"
 	session_duration = "{{.App.SessionDuration}}"
}
`

func init() {
	rootCmd.AddCommand(accessApplicationCmd)
}

var accessApplicationCmd = &cobra.Command{
	Use:   "access_application",
	Short: "Import Access Application data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("Importing Access Application data")

		for _, zone := range zones {
			log.Printf("[DEBUG] Processing zone: ID %s, Name %s", zone.ID, zone.Name)

			accessApplications, _, err := api.AccessApplications(zone.ID, cloudflare.PaginationOptions{
				Page:    1,
				PerPage: 1000,
			})

			if err != nil {
				if strings.Contains(err.Error(), "HTTP status 403") {
					log.Printf("[INFO] insufficient permissions accessing Zone ID %s\n", zone.ID)
					continue
				}

				fmt.Println(err)
				os.Exit(1)
			}

			for _, app := range accessApplications {

				accessApplicationParse(app, zone)
			}

		}
	},
}

func accessApplicationParse(app cloudflare.AccessApplication, zone cloudflare.Zone) {
	tmpl := template.Must(template.New("access_rule").Funcs(templateFuncMap).Parse(accessApplicationTemplate))
	tmpl.Execute(os.Stdout,
		struct {
			App  cloudflare.AccessApplication
			Zone cloudflare.Zone
		}{
			App:  app,
			Zone: zone,
		})
}
