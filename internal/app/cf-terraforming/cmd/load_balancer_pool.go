package cmd

import (
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go"

	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const loadBalancerPoolTemplate = `
resource "cloudflare_load_balancer_pool" "load_balancer_pool_{{.LBP.ID}}" {
    name = "{{.LBP.Name}}"
{{- if .LBP.Origins}}
    {{- range .LBP.Origins}}
    origins {
        name = "{{.Name}}"
        address = "{{.Address}}"
        weight = {{.Weight}}
        enabled = {{.Enabled}}
    }
    {{- end }}
{{- end }}
{{- if .LBP.Description }}
    description = "{{.LBP.Description}}"
{{- end }}
{{- if .LBP.Enabled }}
    enabled = {{.LBP.Enabled}}
{{- end }}
{{- if .LBP.MinimumOrigins }}
    minimum_origins = {{.LBP.MinimumOrigins}}
{{- end }}
{{- if .LBP.Monitor }}
    monitor = "{{.LBP.Monitor}}"
{{- end }}
{{- if .LBP.NotificationEmail }}
    notification_email = "{{.LBP.NotificationEmail}}"
{{- end }}
}
`

func init() {
	rootCmd.AddCommand(loadBalancerPoolCmd)
}

var loadBalancerPoolCmd = &cobra.Command{
	Use:   "load_balancer_pool",
	Short: "Import a load balancer pool into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		loadBalancerPools, err := api.ListLoadBalancerPools()

		if err != nil {
			log.Error(err)
			return
		}

		if len(loadBalancerPools) > 0 {
			for _, lbp := range loadBalancerPools {

				log.WithFields(logrus.Fields{
					"ID":          lbp.ID,
					"Description": lbp.Description,
				}).Debug("Processing load balancer pool")

				if tfstate {
					// TODO: Implement state dump
				} else {
					loadBalancerPoolParse(lbp)
				}
			}
		}
	},
}

func loadBalancerPoolParse(lbp cloudflare.LoadBalancerPool) {
	tmpl := template.Must(template.New("load_balancer_pool").Funcs(templateFuncMap).Parse(loadBalancerPoolTemplate))
	err := tmpl.Execute(os.Stdout,
		struct {
			LBP cloudflare.LoadBalancerPool
		}{
			LBP: lbp,
		})
	if err != nil {
		log.Error(err)
	}
}
