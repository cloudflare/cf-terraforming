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
            {{- if isSlice $v }}
                {{- $k }} = [ {{ range $v }}"{{.}}", {{ end }} ]
            {{- else if isMap $v }}
                {{ $k }} {
                {{- range $k1, $v1 := $v}}
                    {{- if isSlice $v1 }}
                        {{ $k1 }} = [ {{ range $v1 }}"{{.}}", {{ end }} ]
                    {{- else }}
                        {{$k1}} = {{ quoteIfString $v1 -}}
                    {{- end }}
                {{- end }}
                }
            {{- else }}
                {{$k}} = {{ quoteIfString $v -}}
            {{- end }}
        {{- end }}
        }
    {{ else if isSlice .Value}}
        {{- .ID }} = [ {{ range .Value }}"{{.}}", {{ end }} ]
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
				log.Error(err)
				return
			}

			for _, rule := range pageRules {

				log.WithFields(logrus.Fields{
					"ID":       rule.ID,
					"Targets":  rule.Targets,
					"Priority": rule.Priority,
					"Status":   rule.Status,
				}).Debug("Processing page rule")

				actionsProcessed := make([]cloudflare.PageRuleAction, 0)
				for _, action := range rule.Actions {

					if action.ID == "cache_key_fields" {
						for fieldID, fieldValue := range action.Value.(map[string]interface{}) {
							if fieldID == "query_string" {
								fieldMap := fieldValue.(map[string]interface{})

								if fieldMap["exclude"] == "*" {
									fieldMap["exclude"] = []interface{}{}
									fieldMap["ignore"] = true
								}

								if fieldMap["include"] == "*" {
									fieldMap["include"] = []interface{}{}
									fieldMap["ignore"] = false
								}
							}
						}

						actionsProcessed = append(actionsProcessed, action)

					} else if action.ID == "cache_ttl_by_status" {
						for statusCodes, statusTTL := range action.Value.(map[string]interface{}) {
							entry := map[string]interface{}{"codes": statusCodes}

							switch value := statusTTL.(type) {
							case float64:
								entry["ttl"] = int32(value)
							case string:
								switch value {
								case "no-cache":
									entry["ttl"] = 0
								case "no-store":
									entry["ttl"] = -1
								}
							}

							actionsProcessed = append(actionsProcessed, cloudflare.PageRuleAction{
								ID:    action.ID,
								Value: entry,
							})
						}

					} else {
						actionsProcessed = append(actionsProcessed, action)
					}

					rule.Actions = actionsProcessed
				}

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
	err := tmpl.Execute(os.Stdout,
		struct {
			Rule cloudflare.PageRule
			Zone cloudflare.Zone
		}{
			Rule: rule,
			Zone: zone,
		})
	if err != nil {
		log.Error(err)
	}
}
