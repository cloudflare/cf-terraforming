package cmd

import (
	"os"
	"strconv"

	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const firewallRuleTemplate = `
resource "cloudflare_firewall_rule" "firewall_rule_{{.FirewallRule.ID}}" {
  zone_id = "{{.Zone.ID}}"
  description = "{{.FirewallRule.Description}}"
  filter_id = "{{.FirewallRule.Filter.ID}}"
  action = "{{.FirewallRule.Action}}"
  {{- if .FirewallRule.Priority}}
  priority = {{.FirewallRule.Priority}}
  {{- end }}
  {{- if .FirewallRule.Paused}}
  paused = {{.FirewallRule.Paused}}
  {{- end }}
  {{- if .FirewallRule.Products}}
  products = [{{range .FirewallRule.Products}}"{{.}}",{{end}}]
  {{- end }}
}
`

type FirewallRuleAttributes struct {
	ID          string   `json:"id"`
	Action      string   `json:"action"`
	FilterID    string   `json:"filter_id"`
	Priority    string   `json:"priority"`
	ZoneID      string   `json:"zone_id"`
	Description string   `json:"description"`
	Paused      string   `json:"paused"`
	Products    []string `json:"products,omitempty"`
}

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
				log.Error(err)
				return
			}

			for _, r := range firewallRules {

				log.WithFields(logrus.Fields{
					"ID":          r.ID,
					"Description": r.Description,
				}).Debug("Processing firewall rule")

				if tfstate {
					r := firewallRuleResourceStateBuild(zone, r)
					resourcesMap["cloudflare_firewall_rule.firewall_rule_"+r.Primary.Id] = r
				} else {
					firewallRuleParse(zone, r)
				}
			}
		}
	},
}

func firewallRuleParse(zone cloudflare.Zone, firewallRule cloudflare.FirewallRule) {
	tmpl := template.Must(template.New("firewall_rule").Funcs(templateFuncMap).Parse(firewallRuleTemplate))
	err := tmpl.Execute(os.Stdout,
		struct {
			Zone         cloudflare.Zone
			FirewallRule cloudflare.FirewallRule
		}{
			Zone:         zone,
			FirewallRule: firewallRule,
		})
	if err != nil {
		log.Error(err)
	}
}

func firewallRuleResourceStateBuild(zone cloudflare.Zone, rule cloudflare.FirewallRule) Resource {
	r := Resource{
		Primary: Primary{
			Id: rule.ID,
			Attributes: FirewallRuleAttributes{
				ID:          rule.ID,
				Action:      rule.Action,
				FilterID:    rule.Filter.ID,
				ZoneID:      zone.ID,
				Description: rule.Description,
				Paused:      strconv.FormatBool(rule.Paused),
				Products:    rule.Products,
			},
			Meta:    make(map[string]string),
			Tainted: false,
		},
		DependsOn: []string{},
		Deposed:   []string{},
		Provider:  "provider.cloudflare",
		Type:      "cloudflare_firewall_rule",
	}

	firewallRuleAttributes := r.Primary.Attributes.(FirewallRuleAttributes)
	if rule.Priority != nil {
		firewallRuleAttributes.Priority = strconv.FormatFloat(rule.Priority.(float64), 'g', -1, 32)
	}
	r.Primary.Attributes = firewallRuleAttributes

	return r
}
