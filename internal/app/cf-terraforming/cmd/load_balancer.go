package cmd

import (
	"fmt"
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go"

	"text/template"

	"github.com/spf13/cobra"
)

const loadBalancerTemplate = `
resource "cloudflare_load_balancer" "{{.LB.ID}}" {
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
		// Loop through all zones in account and fetch routes for each zone
		for _, zone := range zones {
			loadBalancers, err := api.ListLoadBalancers(zone.ID)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if len(loadBalancers) > 0 {
				for _, lb := range loadBalancers {
					loadBalancerParse(lb, zone)
				}
			}
		}
	},
}

func loadBalancerParse(lb cloudflare.LoadBalancer, zone cloudflare.Zone) {
	tmpl := template.Must(template.New("script").Funcs(templateFuncMap).Parse(loadBalancerTemplate))
	tmpl.Execute(os.Stdout,
		struct {
			LB   cloudflare.LoadBalancer
			Zone cloudflare.Zone
		}{
			LB:   lb,
			Zone: zone,
		})
}
