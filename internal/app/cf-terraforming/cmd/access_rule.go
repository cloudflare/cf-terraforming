package cmd

import (
	"os"
	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const accessRuleTemplate = `
resource "cloudflare_access_rule" "access_rule_{{.AccessRule.ID}}" {
  notes = "{{.AccessRule.Notes}}"
  mode = "{{.AccessRule.Mode}}"
  configuration = {
    target = "{{.AccessRule.Configuration.Target}}"
    value = "{{.AccessRule.Configuration.Value}}"
  }
  {{- if eq .AccessRule.Scope.Type "zone"}}
  zone_id = "{{.Zone.ID}}"
  {{- end }}
}
`

type AccessRuleAttributes struct {
	ID               string `json:"id"`
	Notes            string `json:"notes"`
	Mode             string `json:"mode"`
	ConfigurationNum string `json:"configuration.%"`
	Target           string `json:"configuration.target"`
	Value            string `json:"configuration.value"`
	ZoneID           string `json:"zone_id,omitempty"`
}

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
					log.Error(err)
				}

				totalPages = accessRules.TotalPages

				for _, r := range accessRules.Result {

					log.WithFields(logrus.Fields{
						"Rule ID":       r.ID,
						"Notes":         r.Notes,
						"Configuration": r.Configuration,
						"Scope":         r.Scope,
					}).Debug("Processing Access rule")

					if tfstate {
						state := accessRuleResourceStateBuild(zone, r)
						resourcesMap["cloudflare_access_rule.access_rule_"+r.ID] = state
					} else {
						accessRuleParse(zone, r)
					}
				}
			}
		}
	},
}

func accessRuleParse(zone cloudflare.Zone, accessRule cloudflare.AccessRule) {
	tmpl := template.Must(template.New("access_rule").Funcs(templateFuncMap).Parse(accessRuleTemplate))
	err := tmpl.Execute(os.Stdout,
		struct {
			Zone       cloudflare.Zone
			AccessRule cloudflare.AccessRule
		}{
			Zone:       zone,
			AccessRule: accessRule,
		})
	if err != nil {
		log.Error(err)
	}
}

func accessRuleResourceStateBuild(zone cloudflare.Zone, rule cloudflare.AccessRule) Resource {
	r := Resource{
		Primary: Primary{
			Id: rule.ID,
			Attributes: AccessRuleAttributes{
				ID:               rule.ID,
				Notes:            rule.Notes,
				Mode:             rule.Mode,
				ConfigurationNum: "2",
				Target:           rule.Configuration.Target,
				Value:            rule.Configuration.Value,
			},
			Meta:    make(map[string]string),
			Tainted: false,
		},
		DependsOn: []string{},
		Deposed:   []string{},
		Provider:  "provider.cloudflare",
		Type:      "cloudflare_access_rule",
	}

	attributes := r.Primary.Attributes.(AccessRuleAttributes)

	if rule.Scope.Type == "zone" {
		attributes.ZoneID = zone.ID
	}

	r.Primary.Attributes = attributes

	return r
}
