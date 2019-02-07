package cmd

import (
	"fmt"
	"log"
	"os"

	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/spf13/cobra"
)

const zoneTemplate = `
resource "cloudflare_zone" "{{replace .Zone.Name "." "_"}}" {
    zone = "{{.Zone.Name}}"
{{ if .Zone.Paused}}    paused = "true"{{end}}
    plan = "{{.ZonePlan}}"
}

`

// we enforce the use of the Cloudflare API 'legacy_id' field until the mapping of plan is fixed in cloudflare-go
const (
	planIDFree       = "free"
	planIDPro        = "pro"
	planIDBusiness   = "business"
	planIDEnterprise = "enterprise"
)

// we keep a private map and we will have a function to check and validate the descriptive name from the RatePlan API with the legacy_id
var idForName = map[string]string{
	"Free Website":       planIDFree,
	"Pro Website":        planIDPro,
	"Business Website":   planIDBusiness,
	"Enterprise Website": planIDEnterprise,
}

func init() {
	rootCmd.AddCommand(zoneCmd)
}

var zoneCmd = &cobra.Command{
	Use:   "zone",
	Short: "Import zone data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("Importing zones' data")

		for _, zone := range zones {
			zoneDetails, err := api.ZoneDetails(zone.ID)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			log.Printf("[DEBUG] Processing zone: ID %s, Name %s", zoneDetails.ID, zoneDetails.Name)

			zoneParse(zone)
		}
	},
}

func zoneParse(zone cloudflare.Zone) {
	tmpl := template.Must(template.New("zone").Funcs(templateFuncMap).Parse(zoneTemplate))
	tmpl.Execute(os.Stdout,
		struct {
			Zone     cloudflare.Zone
			ZonePlan string
		}{
			Zone:     zone,
			ZonePlan: idForName[zone.Plan.Name],
		})
}
