package cmd

import (
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go"

	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const loadBalancerMonitorTemplate = `
resource "cloudflare_load_balancer_monitor" "{{.LBM.ID}}" {
    expected_body = "{{.LBM.ExpectedBody}}"
    expected_codes = "{{.LBM.ExpectedCodes}}"
    method = "{{.LBM.Method}}"
    timeout = {{.LBM.Timeout}}
    path = "{{.LBM.Path}}"
    interval = {{.LBM.Interval}}
    retries = {{.LBM.Retries}}
    description = "{{.LBM.Description}}"
    {{if isMap .LBM.Header}}
    header {
    {{range $k, $v := .LBM.Header}}
        {{$k}} = {{ quoteIfString $v }}
    {{end}}
    }
    {{end}}
}
`

func init() {
	rootCmd.AddCommand(loadBalancerMonitorCmd)
}

var loadBalancerMonitorCmd = &cobra.Command{
	Use:   "load_balancer_monitor",
	Short: "Import a load balancer monitor into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		loadBalancerMonitors, err := api.ListLoadBalancerMonitors()

		if err != nil {
			log.Debug(err)
		}

		if len(loadBalancerMonitors) > 0 {
			for _, lbm := range loadBalancerMonitors {

				log.WithFields(logrus.Fields{
					"ID":          lbm.ID,
					"Description": lbm.Description,
				}).Debug("Processing load balancer monitor")

				if tfstate {
					// TODO: Implement state dump
				} else {
					loadBalancerMonitorParse(lbm)
				}
			}
		}
	},
}

func loadBalancerMonitorParse(lbm cloudflare.LoadBalancerMonitor) {
	tmpl := template.Must(template.New("load_balancer_monitor").Funcs(templateFuncMap).Parse(loadBalancerMonitorTemplate))
	tmpl.Execute(os.Stdout,
		struct {
			LBM cloudflare.LoadBalancerMonitor
		}{
			LBM: lbm,
		})
}
