package cmd

import (
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go"

	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const workerRouteTemplate = `
resource "cloudflare_worker_route" "{{.Route.ID}}" {
    zone = "{{.Zone.Name}}"
    pattern = "{{.Route.Pattern}}"
{{if .MultiScript }}
	script_name = "${cloudflare_worker_script.{{.Route.Script}}}"
{{else}}
    enabled = "{{.Route.Enabled}}"
{{end}}
}
`

func init() {
	rootCmd.AddCommand(workerRouteCmd)
}

var workerRouteCmd = &cobra.Command{
	Use:   "worker_route",
	Short: "Import a worker route into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Importing worker route data")
		// Loop through all zones in account and fetch routes for each zone
		for _, zone := range zones {
			workerRoutesResponse, err := api.ListWorkerRoutes(zone.ID)

			if err != nil {
				log.Debug(err)
				return
			}

			if workerRoutesResponse.Success == true {
				for _, route := range workerRoutesResponse.Routes {

					log.WithFields(logrus.Fields{
						"ID":      route.ID,
						"Pattern": route.Pattern,
					}).Debug("Processing woker route")
					// worker_route is rendered differently for multi-script (enterprise) accounts
					// and non-enterprise accounts
					workerRouteParse(zone, route, api.OrganizationID != "")
				}
			}
		}

	},
}

func workerRouteParse(zone cloudflare.Zone, route cloudflare.WorkerRoute, multiScript bool) {
	tmpl := template.Must(template.New("script").Funcs(templateFuncMap).Parse(workerRouteTemplate))
	tmpl.Execute(os.Stdout,
		struct {
			Zone        cloudflare.Zone
			Route       cloudflare.WorkerRoute
			MultiScript bool
		}{
			Zone:        zone,
			Route:       route,
			MultiScript: multiScript,
		})
}
