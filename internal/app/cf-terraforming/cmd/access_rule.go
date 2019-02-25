package cmd

import (
	"os"
	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const accessRuleTemplate = `
resource "cloudflare_access_rule" "{{.AccessRule.ID}}" {
  notes = "{{.AccessRule.Notes}}"
  mode = "{{.AccessRule.Mode}}"
  configuration {
    target = "{{.AccessRule.Configuration.Target}}"
    value = "{{.AccessRule.Configuration.Value}}"
  }
  {{if eq .AccessRule.Scope.Type "zone"}}zone_id = "{{.Zone.ID}}"{{end}}
}
`

func init() {
	rootCmd.AddCommand(accessRuleCmd)
}

var accessRuleCmd = &cobra.Command{
	Use:   "access_rule",
	Short: "Import Access Rule data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Importing Access Rule data")

		for _, zone := range zones {

			log.WithFields(logrus.Fields{
				"ID":   zone.ID,
				"Name": zone.Name,
			}).Debug("Processing zone")

			totalPages := 999

			for page := 1; page <= totalPages; page++ {
				accessRules, err := api.ListZoneAccessRules(zone.ID, cloudflare.AccessRule{}, page)

				if err != nil {
					log.Debug(err)
				}

				totalPages = accessRules.TotalPages

				for _, r := range accessRules.Result {

					log.WithFields(logrus.Fields{
						"Rule ID":       r.ID,
						"Notes":         r.Notes,
						"Configuration": r.Configuration,
						"Scope":         r.Scope,
					}).Debug("Processing Access rule")

					accessRuleParse(zone, r)
				}
			}
		}
	},
}

func accessRuleParse(zone cloudflare.Zone, accessRule cloudflare.AccessRule) {
	tmpl := template.Must(template.New("access_rule").Funcs(templateFuncMap).Parse(accessRuleTemplate))
	tmpl.Execute(os.Stdout,
		struct {
			Zone       cloudflare.Zone
			AccessRule cloudflare.AccessRule
		}{
			Zone:       zone,
			AccessRule: accessRule,
		})
}
