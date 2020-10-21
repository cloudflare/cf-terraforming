package cmd

import (
	"os"
	"sort"
	"strings"

	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const zoneSettingOverrideTemplate = `
resource "cloudflare_zone_settings_override" "zone_settings_override_{{.Zone.ID}}" {
	zone_id = "{{.Zone.ID -}}"
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
		{{ else if isSlice .Value}}
			{{- .ID }} = [ {{ range .Value }}"{{.}}", {{ end }} ]
		{{ else }}
			{{- .ID}} = {{ quoteIfString .Value -}}
		{{- end}}
		{{- end}}
	}
}
`

var zoneSettingsToRemove = []string{
	"advanced_ddos", "edge_cache_ttl",
}

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
				log.Error(err)
				return
			}

			if settingsResponse.Success {
				i := 0

				for _, s := range settingsResponse.Result {
					if contains(zoneSettingsToRemove, s.ID) {
						continue
					}

					// Remap the 0rtt zone setting to zero_rtt
					if s.ID == "0rtt" {
						s.ID = "zero_rtt"
					}

					settingsResponse.Result[i] = s
					i++
				}

				settingsResponse.Result = settingsResponse.Result[:i]
			}

			sort.Slice(settingsResponse.Result, func(i, j int) bool {
				return strings.Compare(settingsResponse.Result[i].ID, settingsResponse.Result[j].ID) <= 0
			})

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
	err := tmpl.Execute(os.Stdout,
		struct {
			Settings []cloudflare.ZoneSetting
			Zone     cloudflare.Zone
		}{
			Settings: s,
			Zone:     zone,
		})
	if err != nil {
		log.Error(err)
	}
}
