package cmd

import (
	"os"

	"strings"
	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const accessPolicyTemplate = `
resource "cloudflare_access_policy" "{{.Policy.ID}}" {
    application_id = "{{.App.ID}}"
    zone_id = "{{.Zone.ID}}"
    name = "{{.Policy.Name}}"
    precedence = "{{.Policy.Precedence}}"
    decision = "{{.Policy.Decision}}"

{{if .Policy.Include }}
    include = {
{{range $k, $v := .Policy.Include }}
    {{if isMap $v }}
        {{- range $k, $v := $v }}
            {{ $k }} =  {{if isMap $v }} [{{range $v}}"{{.}}",{{end}}]  {{else}} "{{ $v }}" {{end}}
        {{end -}}
    {{end}}
{{end}}
    }
{{end}}

{{if .Policy.Exclude }}
    exclude = {
{{range $k, $v := .Policy.Exclude }}
    {{if isMap $v }}
        {{- range $k, $v := $v }}
            {{ $k }} =  {{if isMap $v }} [{{range $v}}"{{.}}",{{end}}]  {{else}} "{{ $v }}" {{end}}
        {{end -}}
    {{end}}
{{end}}
    }
{{end}}

{{if .Policy.Require }}
    require = {
{{range $k, $v := .Policy.Require }}
    {{if isMap $v }}
        {{- range $k, $v := $v }}
            {{ $k }} =  {{if isMap $v }} [{{range $v}}"{{.}}",{{end}}]  {{else}} "{{ $v }}" {{end}}
        {{end -}}
    {{end}}
{{end}}
    }
{{end}}

}
`

func init() {
	rootCmd.AddCommand(accessPolicyCmd)
}

var accessPolicyCmd = &cobra.Command{
	Use:   "access_policy",
	Short: "Import Access Policy data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Importing Access Policy data")

		for _, zone := range zones {
			log.WithFields(logrus.Fields{
				"ID":   zone.ID,
				"Name": zone.Name,
			}).Debug("Processing zone")

			accessApplications, _, appFetchErr := api.AccessApplications(zone.ID, cloudflare.PaginationOptions{
				Page:    1,
				PerPage: 1000,
			})

			if appFetchErr != nil {
				log.Debug(appFetchErr)
				return
			}

			for _, app := range accessApplications {

				accessPolicies, _, err := api.AccessPolicies(zone.ID, app.ID, cloudflare.PaginationOptions{
					Page:    1,
					PerPage: 1000,
				})

				if err != nil {
					if strings.Contains(err.Error(), "HTTP status 403") {
						log.WithFields(logrus.Fields{
							"ID": zone.ID,
						}).Debug("Insufficient permissions for accessing zone")
						continue
					}
					log.Debug(err)
				}

				for _, policy := range accessPolicies {

					accessPolicyParse(app, policy, zone)
				}
			}

		}
	},
}

func accessPolicyParse(app cloudflare.AccessApplication, policy cloudflare.AccessPolicy, zone cloudflare.Zone) {
	tmpl := template.Must(template.New("access_policy").Funcs(templateFuncMap).Parse(accessPolicyTemplate))
	tmpl.Execute(os.Stdout,
		struct {
			App    cloudflare.AccessApplication
			Policy cloudflare.AccessPolicy
			Zone   cloudflare.Zone
		}{
			App:    app,
			Policy: policy,
			Zone:   zone,
		})
}
