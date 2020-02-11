package cmd

import (
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go"

	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const loadBalancerTemplate = `
resource "cloudflare_load_balancer" "load_balancer_{{.LB.ID}}" {
    zone = "{{.Zone.Name}}"
    name = "{{.LB.Name}}"
    fallback_pool_id = "{{.LB.FallbackPool}}"
{{if .LB.DefaultPools}}
    default_pool_ids = [{{range .LB.DefaultPools}}"{{.}}",{{end}}]
{{end}}
{{if .LB.Description}}
    description = "{{.LB.Description}}"
{{end}}
{{/* TTL conflicts with Proxied setting and cannot be set for a Proxied LB */}}
{{/* See: https://www.terraform.io/docs/providers/cloudflare/r/load_balancer.html */}}
{{if and (ne .LB.TTL 0) (eq .LB.Proxied true) }}
    ttl = {{.LB.TTL}}
{{end}}
{{if .LB.SteeringPolicy}}
    steering_policy = "{{.LB.SteeringPolicy}}"
{{end}}
{{if .LB.Proxied}}
    proxied = {{.LB.Proxied}}
{{end}}
{{if .LB.RegionPools}}
    {{range $region, $regIDs := .LB.RegionPools}}
        region_pools {
            region = "{{ $region }}"
            pool_ids = [{{ range $regIDs }} "{{.}}", {{end}}]
        }
    {{end}}
{{end}}
{{if .LB.PopPools}}
    {{range $pop, $popIDs := .LB.PopPools}}
        pop_pools {
            pop = "{{ $pop }}"
            pool_ids = [{{range $popID := $popIDs }} "{{.}}", {{end}}]
        }
    {{end}}
{{end}}
{{if .LB.Persistence}}
    session_affinity = "{{.LB.Persistence}}"
{{end}}
}
`

func init() {
	rootCmd.AddCommand(loadBalancerCmd)
}

var loadBalancerCmd = &cobra.Command{
	Use:   "load_balancer",
	Short: "Import a load balancer into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Importing Load Balancer data")
		// Loop through all zones in account and fetch routes for each zone
		for _, zone := range zones {
			loadBalancers, err := api.ListLoadBalancers(zone.ID)

			log.WithFields(logrus.Fields{
				"ID":   zone.ID,
				"Name": zone.Name,
			})

			if err != nil {
				log.Debug(err)
				return
			}

			if len(loadBalancers) > 0 {
				for _, lb := range loadBalancers {

					log.WithFields(logrus.Fields{
						"ID":           lb.ID,
						"Description":  lb.Description,
						"FallbackPool": lb.FallbackPool,
						"DefaultPools": lb.DefaultPools,
					}).Debug("Processing load balancer")

					if tfstate {
						// TODO: Implement state dump
					} else {
						loadBalancerParse(lb, zone)
					}
				}
			}
		}
	},
}

func loadBalancerParse(lb cloudflare.LoadBalancer, zone cloudflare.Zone) {
	tmpl := template.Must(template.New("script").Funcs(templateFuncMap).Parse(loadBalancerTemplate))
	err := tmpl.Execute(os.Stdout,
		struct {
			LB   cloudflare.LoadBalancer
			Zone cloudflare.Zone
		}{
			LB:   lb,
			Zone: zone,
		})
	if err != nil {
		log.Error(err)
	}
}
