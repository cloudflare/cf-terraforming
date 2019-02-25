package cmd

import (
	"os"

	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const wafRuleTemplate = `
resource "cloudflare_waf_rule" "{{replace .Zone.Name "." "_"}}_{{.Rule.ID}}" {
    rule_id = "{{.Rule.ID}}"
    zone = "{{.Zone.Name}}"
    mode = "{{.Rule.Mode}}"
}
`

func init() {
	rootCmd.AddCommand(wafRuleCmd)
}

var wafRuleCmd = &cobra.Command{
	Use:   "waf_rule",
	Short: "Import WAF Rule data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Importing WAF Rule data")

		for _, zone := range zones {

			log.WithFields(logrus.Fields{
				"ID":   zone.ID,
				"Name": zone.Name,
			}).Debug("Processing zone")

			wafPackages, err := api.ListWAFPackages(zone.ID)

			if err != nil {
				log.Debug(err)
				return
			}

			for _, wafPackage := range wafPackages {

				log.WithFields(logrus.Fields{
					"ID":   wafPackage.ID,
					"Name": wafPackage.Name,
				}).Debug("Processing WAF package")

				log.Debug("Fetching WAF rules")

				wafRules, err := api.ListWAFRules(zone.ID, wafPackage.ID)

				if err != nil {
					log.Debug(err)
					return
				}

				for _, rule := range wafRules {

					log.WithFields(logrus.Fields{
						"ID": rule.ID,
					}).Debug("Processing WAF rule")

					wafRuleParse(zone, wafPackage, rule)
				}
			}
		}
	},
}

func wafRuleParse(zone cloudflare.Zone, wafPackage cloudflare.WAFPackage, wafRule cloudflare.WAFRule) {
	tmpl := template.Must(template.New("waf_rule").Funcs(templateFuncMap).Parse(wafRuleTemplate))
	tmpl.Execute(os.Stdout,
		struct {
			Zone    cloudflare.Zone
			Package cloudflare.WAFPackage
			Rule    cloudflare.WAFRule
		}{
			Zone:    zone,
			Package: wafPackage,
			Rule:    wafRule,
		})
}
