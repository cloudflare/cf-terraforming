package cmd

import (
	"os"

	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const pageRuleTemplate = `
resource "cloudflare_page_rule" "page_rule_{{.Rule.ID}}" {
    zone_id = "{{.Zone.ID}}"
{{ range .Rule.Targets}}
    target = "{{.Constraint.Value }}"
{{ end }}
    priority = {{ quoteIfString .Rule.Priority }}
    status = "{{.Rule.Status}}"
    actions {
    {{- range .Rule.Actions}}
    {{- if isMap .Value}}
        {{.ID}} {
        {{- range $k, $v := .Value }}
            {{$k}} = {{ quoteIfString $v -}}
        {{end }}
        }
    {{else}}
        {{.ID}} = {{ quoteIfString .Value }}
    {{end -}}
    {{end }}
    }
}
`

func init() {
	rootCmd.AddCommand(pageRuleCmd)
}

var pageRuleCmd = &cobra.Command{
	Use:   "page_rule",
	Short: "Import Page Rule data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Importing Page Rule data")

		for _, zone := range zones {

			log.WithFields(logrus.Fields{
				"ID":   zone.ID,
				"Name": zone.Name,
			}).Debug("Processing zone")

			pageRules, err := api.ListPageRules(zone.ID)

			if err != nil {
				log.Debug(err)
				return
			}

			for _, rule := range pageRules {

				log.WithFields(logrus.Fields{
					"ID":       rule.ID,
					"Targets":  rule.Targets,
					"Priority": rule.Priority,
					"Status":   rule.Status,
				}).Debug("Processing page rule")

				if tfstate {
					// TODO: Implement state dump
				} else {
					pageRuleParse(rule, zone)
				}
			}
		}
	},
}

func pageRuleParse(rule cloudflare.PageRule, zone cloudflare.Zone) {
	tmpl := template.Must(template.New("page_rule").Funcs(templateFuncMap).Parse(pageRuleTemplate))
	tmpl.Execute(os.Stdout,
		struct {
			Rule cloudflare.PageRule
			Zone cloudflare.Zone
		}{
			Rule: rule,
			Zone: zone,
		})
}
