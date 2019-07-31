package cmd

import (
	"os"
	"strconv"

	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const workerRouteTemplate = `
resource "cloudflare_worker_route" "worker_route_{{.Route.ID}}" {
    zone = "{{.Zone.Name}}"
    pattern = "{{.Route.Pattern}}"
{{if .MultiScript }}
	script_name = "${cloudflare_worker_script.{{.Route.Script}}}"
{{else}}
    enabled = "{{.Route.Enabled}}"
{{end}}
}
`

type WorkerRouteAttributes struct {
	Enabled     string `json:"enabled"`
	Id          string `json:"id"`
	MultiScript string `json:"multi_script"`
	Pattern     string `json:"pattern"`
	Zone        string `json:"zone"`
	ZoneId      string `json:"zone_id"`
}

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

					if tfstate {
						r := workerResourceStateBuild(zone, route, api.OrganizationID != "")

						resourcesMap["cloudflare_worker_route.worker_route_"+route.ID] = r

					} else {
						// worker_route is rendered differently for multi-script (enterprise) accounts
						// and non-enterprise accounts
						workerRouteParse(zone, route, api.OrganizationID != "")
					}
				}
			}
		}
	},
}

func workerResourceStateBuild(zone cloudflare.Zone, route cloudflare.WorkerRoute, multiScript bool) Resource {

	r := Resource{
		Primary: Primary{
			Id: route.ID,
			Attributes: WorkerRouteAttributes{
				Enabled:     strconv.FormatBool(route.Enabled),
				Id:          route.ID,
				MultiScript: strconv.FormatBool(multiScript),
				Pattern:     route.Pattern,
				Zone:        zone.Name,
				ZoneId:      zone.ID,
			},
			Meta:    make(map[string]string),
			Tainted: false,
		},
		DependsOn: []string{},
		Deposed:   []string{},
		Provider:  "provider.cloudflare",
		Type:      "cloudflare_worker_route",
	}

	return r
}

func workerRouteParse(zone cloudflare.Zone, route cloudflare.WorkerRoute, multiScript bool) {

	tmpl := template.Must(template.New("script").Funcs(templateFuncMap).Parse(workerRouteTemplate))
	tmpl.Execute(os.Stdout,
		struct {
			Zone        cloudflare.Zone
			Route       cloudflare.WorkerRoute
			MultiScript bool
			StateHeader string
		}{
			Zone:        zone,
			Route:       route,
			MultiScript: multiScript,
		})
}
