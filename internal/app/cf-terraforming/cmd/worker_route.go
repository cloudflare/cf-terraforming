package cmd

import (
	"fmt"
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go"

	"text/template"

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
		// Loop through all zones in account and fetch routes for each zone
		for _, zone := range zones {
			workerRoutesResponse, err := api.ListWorkerRoutes(zone.ID)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if workerRoutesResponse.Success == true {
				for _, route := range workerRoutesResponse.Routes {
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
