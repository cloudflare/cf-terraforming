package cmd

import (
	"fmt"
	"log"
	"os"

	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/spf13/cobra"
)

const rateLimitTemplate = `
resource "cloudflare_rate_limit" "{{replace .Zone.Name "." "_"}}_{{.RateLimit.ID}}" {
  zone = "{{.Zone.Name}}"
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
    timeout = {{.RateLimit.Action.Timeout}}
	{{if .RateLimit.Action.Response}}
    response {
      content_type = "{{.RateLimit.Action.Response.ContentType}}"
      body = "{{js .RateLimit.Action.Response.Body}}"
    }
	{{end}}
  }
  {{if .RateLimit.Correlate.By}}
  correlate {
    by = "{{.RateLimit.Correlate.By}}"
  }
  {{end}}
  disabled = {{.RateLimit.Disabled}}
  description = "{{.RateLimit.Description}}"
  {{if .RateLimit.Bypass}}
  bypass_url_patterns = [{{range .RateLimit.Bypass.Value}}"{{.}}", {{end}}]
  {{end}}
}
`

func init() {
	rootCmd.AddCommand(rateLimitCmd)
}

var rateLimitCmd = &cobra.Command{
	Use:   "rate_limit",
	Short: "Import Rate Limit data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("Importing Rate Limit data")

		for _, zone := range zones {
			log.Printf("[DEBUG] Processing zone: ID %s, Name %s", zone.ID, zone.Name)

			totalPages := 999

			for page := 1; page <= totalPages; page++ {
				rateLimits, resultInfo, err := api.ListRateLimits(zone.ID, cloudflare.PaginationOptions{
					Page:    page,
					PerPage: 1000,
				})

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				totalPages = resultInfo.TotalPages

				for _, r := range rateLimits {
					log.Printf("[DEBUG] Rate Limit ID %s, Description %s\n", r.ID, r.Description)
					rateLimitParse(zone, r)
				}
			}

		}
	},
}

func rateLimitParse(zone cloudflare.Zone, rateLimit cloudflare.RateLimit) {
	tmpl := template.Must(template.New("rate_limit").Funcs(templateFuncMap).Parse(rateLimitTemplate))
	tmpl.Execute(os.Stdout,
		struct {
			Zone      cloudflare.Zone
			RateLimit cloudflare.RateLimit
		}{
			Zone:      zone,
			RateLimit: rateLimit,
		})
}
