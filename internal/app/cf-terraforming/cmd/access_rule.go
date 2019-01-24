package cmd

import (
	"fmt"
	"log"
	"os"
	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
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
		log.Print("Importing Access Rule data")

		for _, zone := range zones {
			log.Printf("[DEBUG] Processing zone: ID %s, Name %s", zone.ID, zone.Name)

			totalPages := 999

			for page := 1; page <= totalPages; page++ {
				accessRules, err := api.ListZoneAccessRules(zone.ID, cloudflare.AccessRule{}, page)

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				totalPages = accessRules.TotalPages

				for _, r := range accessRules.Result {
					log.Printf("[DEBUG] Access Rule ID %s, Notes %s, Configuration %s, Scope %s\n", r.ID, r.Notes, r.Configuration, r.Scope)
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
