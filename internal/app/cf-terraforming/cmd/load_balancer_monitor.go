package cmd

import (
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go"

	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const loadBalancerMonitorTemplate = `
resource "cloudflare_load_balancer_monitor" "load_balancer_monitor_{{.LBM.ID}}" {
    type = "{{.LBM.Type}}"
    method = "{{.LBM.Method}}"
    timeout = {{.LBM.Timeout}}
    {{- if not (eq .LBM.Port 0) }}
    port = {{.LBM.Port}}
    {{- end }}
    interval = {{.LBM.Interval}}
    retries = {{.LBM.Retries}}
    description = "{{.LBM.Description}}"
    {{- if or (eq .LBM.Type "http") (eq .LBM.Type "https") }}
    expected_body = "{{.LBM.ExpectedBody}}"
    expected_codes = "{{.LBM.ExpectedCodes}}"
    path = "{{.LBM.Path}}"
    {{- if not (isMapEmpty .LBM.Header) }}
    header {
    {{- range $k, $v := .LBM.Header}}
        header = "{{$k}}"
        values = [{{ range $hv := $v }}"{{ $hv }}",{{ end }}]
    {{- end }}
    }
    {{- end }}
    allow_insecure = {{.LBM.AllowInsecure}}
    follow_redirects = {{.LBM.FollowRedirects}}
    #probe_zone = "{{.LBM.ProbeZone}}"
    {{- end }}
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
			log.Error(err)
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
	err := tmpl.Execute(os.Stdout,
		struct {
			LBM cloudflare.LoadBalancerMonitor
		}{
			LBM: lbm,
		})
	if err != nil {
		log.Error(err)
	}
}
