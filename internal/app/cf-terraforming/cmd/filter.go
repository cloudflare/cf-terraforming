package cmd

import (
	"fmt"
	"log"
	"os"
	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/spf13/cobra"
)

const filterTemplate = `
resource "cloudflare_filter" "{{.Filter.ID}}" {
  zone_id = "{{.Zone.ID}}"
  description = "{{.Filter.Description}}"
  expression = "{{js .Filter.Expression}}"
  {{if .Filter.Paused}}paused = {{.Filter.Paused}}{{end}}
  {{if .Filter.Ref}}ref = "{{.Filter.Ref}}"{{end}}
}
`

func init() {
	rootCmd.AddCommand(filterCmd)
}

var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Import Filter data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("Importing Filter data")

		for _, zone := range zones {
			log.Printf("[DEBUG] Processing zone: ID %s, Name %s", zone.ID, zone.Name)

			filters, err := api.Filters(zone.ID, cloudflare.PaginationOptions{
				Page:    1,
				PerPage: 1000,
			})

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for _, r := range filters {
				log.Printf("[DEBUG] Filter ID %s, Expression %s, Description %s\n", r.ID, r.Expression, r.Description)
				filterParse(zone, r)
			}
		}
	},
}

func filterParse(zone cloudflare.Zone, filter cloudflare.Filter) {
	tmpl := template.Must(template.New("filter").Funcs(templateFuncMap).Parse(filterTemplate))
	tmpl.Execute(os.Stdout,
		struct {
			Zone   cloudflare.Zone
			Filter cloudflare.Filter
		}{
			Zone:   zone,
			Filter: filter,
		})
}
