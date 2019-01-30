package cmd

import (
	"fmt"
	"log"
	"os"

	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"

	"github.com/spf13/cobra"
)

const zoneSettingOverrideTemplate = `
resource cloudflare_zone_settings_override {{.Zone.ID}} {
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
		log.Print("Importing Zone Settings data")

		for _, zone := range zones {
			// Fetch all settings for a zone
			settingsResponse, err := api.ZoneSettings(zone.ID)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			zoneSettingsOverrideParse(settingsResponse.Result, zone)
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
