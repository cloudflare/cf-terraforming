package cmd

import (
	"fmt"
	"os"

	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const firewallRuleTemplate = `
resource "cloudflare_firewall_rule" "{{.FirewallRule.ID}}" {
  zone_id = "{{.Zone.ID}}"
  description = "{{.FirewallRule.Description}}"
  filter_id = "{{.FirewallRule.Filter.ID}}"
  action = "{{.FirewallRule.Action}}"
  {{if .FirewallRule.Priority}}priority = {{.FirewallRule.Priority}}{{end}}
  {{if .FirewallRule.Paused}}paused = {{.FirewallRule.Paused}}{{end}}
}
`

func init() {
	rootCmd.AddCommand(firewallRuleCmd)
}

var firewallRuleCmd = &cobra.Command{
	Use:   "firewall_rule",
	Short: "Import Firewall Rule data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Importing Firewall Rule data")

		for _, zone := range zones {

			log.WithFields(logrus.Fields{
				"ID":   zone.ID,
				"Name": zone.Name,
			}).Debug("Processing zone")

			firewallRules, err := api.FirewallRules(zone.ID, cloudflare.PaginationOptions{
				Page:    1,
				PerPage: 1000,
			})

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for _, r := range firewallRules {

				log.WithFields(logrus.Fields{
					"ID":          r.ID,
					"Description": r.Description,
				}).Debug("Processing firewall rule")

				firewallRuleParse(zone, r)
			}
		}
	},
}

func firewallRuleParse(zone cloudflare.Zone, firewallRule cloudflare.FirewallRule) {
	tmpl := template.Must(template.New("firewall_rule").Funcs(templateFuncMap).Parse(firewallRuleTemplate))
	tmpl.Execute(os.Stdout,
		struct {
			Zone         cloudflare.Zone
			FirewallRule cloudflare.FirewallRule
		}{
			Zone:         zone,
			FirewallRule: firewallRule,
		})
}
