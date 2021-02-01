package cmd

import (
	"fmt"
	"os"
	"strings"

	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const rateLimitTemplate = `
resource "cloudflare_rate_limit" "{{replace .Zone.Name "." "_"}}_{{.RateLimit.ID}}" {
  zone_id = "{{.Zone.ID}}"
  threshold = {{.RateLimit.Threshold}}
  period = {{.RateLimit.Period}}
  match {
    request {
      url_pattern = "{{.RateLimit.Match.Request.URLPattern}}"
      schemes = [{{range .RateLimit.Match.Request.Schemes}}"{{.}}", {{end}}]
      methods = [{{range .RateLimit.Match.Request.Methods}}"{{.}}", {{end}}]
    }
    response {
      statuses = [{{range .RateLimit.Match.Response.Statuses}}{{.}}, {{end}}]
      origin_traffic = {{.RateLimit.Match.Response.OriginTraffic}}
    }
  }
  action {
    mode = "{{.RateLimit.Action.Mode}}"
	{{- if or (eq .RateLimit.Action.Mode "simulate") (eq .RateLimit.Action.Mode "ban") }}
    timeout = {{.RateLimit.Action.Timeout}}
	{{- end }}
    {{- if .RateLimit.Action.Response }}
    response {
      content_type = "{{.RateLimit.Action.Response.ContentType}}"
      body = "{{js .RateLimit.Action.Response.Body}}"
    }
    {{- end }}
  }
  {{- if .RateLimit.Correlate }}
  correlate {
    by = "{{.RateLimit.Correlate.By}}"
  }
  {{- end }}
  disabled = {{.RateLimit.Disabled}}
  description = "{{.RateLimit.Description}}"
  {{- if .RateLimit.Bypass }}
  bypass_url_patterns = [{{range .RateLimit.Bypass}}"{{.Value}}", {{end}}]
  {{- end }}
}
`

func init() {
	rootCmd.AddCommand(rateLimitCmd)
}

var rateLimitCmd = &cobra.Command{
	Use:   "rate_limit",
	Short: "Import Rate Limit data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Importing Rate Limit data")

		for _, zone := range zones {

			log.WithFields(logrus.Fields{
				"ID":   zone.ID,
				"Name": zone.Name,
			}).Debug("Processing zone")

			totalPages := 999

			for page := 1; page <= totalPages; page++ {
				rateLimits, resultInfo, err := api.ListRateLimits(zone.ID, cloudflare.PaginationOptions{
					Page:    page,
					PerPage: 1000,
				})

				if err != nil {
					log.Error(err)
					return
				}

				totalPages = resultInfo.TotalPages

				for _, r := range rateLimits {

					log.WithFields(logrus.Fields{
						"ID":          r.ID,
						"Description": r.Description,
					}).Debug("Processing rate limit")

					if tfstate {
						r := rateLimitResourceStateBuild(zone, r)
						resourcesMap["cloudflare_rate_limit."+strings.ReplaceAll(zone.Name, ".", "_")+"_"+r.Primary.Id] = r
					} else {
						rateLimitParse(zone, r)
					}

				}
			}
		}
	},
}

func rateLimitParse(zone cloudflare.Zone, rateLimit cloudflare.RateLimit) {
	tmpl := template.Must(template.New("rate_limit").Funcs(templateFuncMap).Parse(rateLimitTemplate))
	err := tmpl.Execute(os.Stdout,
		struct {
			Zone      cloudflare.Zone
			RateLimit cloudflare.RateLimit
		}{
			Zone:      zone,
			RateLimit: rateLimit,
		})
	if err != nil {
		log.Error(err)
	}
}

func rateLimitResourceStateBuild(zone cloudflare.Zone, rateLimit cloudflare.RateLimit) Resource {
	r := Resource{
		Primary: Primary{
			Id: rateLimit.ID,
			Attributes: map[string]string{
				"id":                                rateLimit.ID,
				"action.#":                          "1",
				"action.0.mode":                     rateLimit.Action.Mode,
				"action.0.timeout":                  fmt.Sprint(rateLimit.Action.Timeout),
				"bypass_url_patterns.#":             fmt.Sprint(len(rateLimit.Bypass)),
				"description":                       rateLimit.Description,
				"disabled":                          fmt.Sprint(rateLimit.Disabled),
				"match.#":                           "1",
				"match.0.request.#":                 "1",
				"match.0.request.0.methods.#":       fmt.Sprint(len(rateLimit.Match.Request.Methods)),
				"match.0.request.0.schemes.#":       fmt.Sprint(len(rateLimit.Match.Request.Schemes)),
				"match.0.request.0.url_pattern":     rateLimit.Match.Request.URLPattern,
				"match.0.response.#":                "1",
				"match.0.response.0.origin_traffic": fmt.Sprint(*rateLimit.Match.Response.OriginTraffic),
				"match.0.response.0.statuses.#":     fmt.Sprint(len(rateLimit.Match.Response.Statuses)),
				"period":                            fmt.Sprint(rateLimit.Period),
				"threshold":                         fmt.Sprint(rateLimit.Threshold),
				"zone":                              zone.Name,
				"zone_id":                           zone.ID,
			},
			Meta:    make(map[string]string),
			Tainted: false,
		},
		DependsOn: []string{},
		Deposed:   []string{},
		Provider:  "provider.cloudflare",
		Type:      "cloudflare_rate_limit",
	}

	attributes := r.Primary.Attributes.(map[string]string)

	if rateLimit.Action.Response != nil {
		attributes["action.0.response.#"] = "1"
		attributes["action.0.response.0.body"] = rateLimit.Action.Response.Body
		attributes["action.0.response.0.content_type"] = rateLimit.Action.Response.ContentType
	} else {
		attributes["action.0.response.#"] = "0"
	}

	for _, bypass := range rateLimit.Bypass {
		hash := hashcode.String(fmt.Sprintf("%s;", bypass.Value))
		attributes[fmt.Sprintf("bypass_url_patterns.%d", hash)] = bypass.Value
	}

	for _, method := range rateLimit.Match.Request.Methods {
		hash := hashcode.String(fmt.Sprintf("%s;", method))
		attributes[fmt.Sprintf("match.0.request.0.methods.%d", hash)] = method
	}

	for _, scheme := range rateLimit.Match.Request.Schemes {
		hash := hashcode.String(fmt.Sprintf("%s;", scheme))
		attributes[fmt.Sprintf("match.0.request.0.schemes.%d", hash)] = scheme
	}

	for _, status := range rateLimit.Match.Response.Statuses {
		hash := hashcode.String(fmt.Sprintf("%s;", fmt.Sprint(status)))
		attributes[fmt.Sprintf("match.0.response.0.statuses.%d", hash)] = fmt.Sprint(status)
	}

	r.Primary.Attributes = attributes

	return r
}
