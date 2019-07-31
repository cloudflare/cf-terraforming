package cmd

import (
	"os"

	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const zoneSettingOverrideTemplate = `
resource "cloudflare_zone_settings_override" "zone_settings_override_{{.Zone.ID}}" {
	name = "{{.Zone.Name -}}"
	settings {
		{{- range .Settings}}
		{{if isMap .Value }}
			{{- .ID }} {
			{{- range $k, $v := .Value}}
			{{if isMap $v }}
				{{- range $k, $v := $v }}
			{{ $k }} = {{ quoteIfString $v -}}
				{{- end}}
				{{else}}
			{{- $k }} = {{ quoteIfString $v -}}
			{{- end}}
		{{- end}}
		}
		{{ else }}
			{{- .ID}} = {{ quoteIfString .Value -}}
		{{- end}}
		{{- end}}
	}
}
`

func init() {
	rootCmd.AddCommand(zoneSettingsOverrideCmd)
}

var zoneSettingsOverrideCmd = &cobra.Command{
	Use:   "zone_settings_override",
	Short: "Import Zone Settings Override data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Importing zone settings data")

		for _, zone := range zones {
			// Fetch all settings for a zone
			settingsResponse, err := api.ZoneSettings(zone.ID)

			if err != nil {
				log.Debug(err)
				return
			}

			log.WithFields(logrus.Fields{
				"Result": settingsResponse.Result,
			}).Debug("Processing zone settings")

			if tfstate {
				// TODO: Implement state dump
			} else {
				zoneSettingsOverrideParse(settingsResponse.Result, zone)
			}
		}
	},
}

func zoneSettingsOverrideParse(s []cloudflare.ZoneSetting, zone cloudflare.Zone) {
	tmpl := template.Must(template.New("script").Funcs(templateFuncMap).Parse(zoneSettingOverrideTemplate))
	tmpl.Execute(os.Stdout,
		struct {
			Settings []cloudflare.ZoneSetting
			Zone     cloudflare.Zone
		}{
			Settings: s,
			Zone:     zone,
		})
}
