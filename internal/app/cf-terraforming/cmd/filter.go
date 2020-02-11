package cmd

import (
	"os"
	"strconv"
	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const filterTemplate = `
resource "cloudflare_filter" "filter_{{.Filter.ID}}" {
  zone_id = "{{.Zone.ID}}"
  description = "{{.Filter.Description}}"
  expression = "{{js .Filter.Expression}}"
  {{if .Filter.Paused}}paused = {{.Filter.Paused}}{{end}}
  {{if .Filter.Ref}}ref = "{{.Filter.Ref}}"{{end}}
}
`

type FilterAttributes struct {
	ID          string `json:"id"`
	ZoneID      string `json:"zone_id"`
	Description string `json:"description"`
	Expression  string `json:"expression"`
	Paused      string `json:"paused"`
	Ref         string `json:"ref"`
}

func init() {
	rootCmd.AddCommand(filterCmd)
}

var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Import Filter data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Importing Filter data")

		for _, zone := range zones {

			log.WithFields(logrus.Fields{
				"ID":   zone.ID,
				"Name": zone.Name,
			}).Debug("Processing zone")

			filters, err := api.Filters(zone.ID, cloudflare.PaginationOptions{
				Page:    1,
				PerPage: 1000,
			})

			if err != nil {
				log.Error(err)
				return
			}

			for _, r := range filters {

				log.WithFields(logrus.Fields{
					"ID":          r.ID,
					"Expression":  r.Expression,
					"Description": r.Description,
				})

				if tfstate {
					r := filterResourceStateBuild(zone, r)
					resourcesMap["cloudflare_filter.filter_"+r.Primary.Id] = r
				} else {
					filterParse(zone, r)
				}
			}
		}
	},
}

func filterParse(zone cloudflare.Zone, filter cloudflare.Filter) {
	tmpl := template.Must(template.New("filter").Funcs(templateFuncMap).Parse(filterTemplate))
	err := tmpl.Execute(os.Stdout,
		struct {
			Zone   cloudflare.Zone
			Filter cloudflare.Filter
		}{
			Zone:   zone,
			Filter: filter,
		})
	if err != nil {
		log.Error(err)
	}
}

func filterResourceStateBuild(zone cloudflare.Zone, filter cloudflare.Filter) Resource {
	r := Resource{
		Primary: Primary{
			Id: filter.ID,
			Attributes: FilterAttributes{
				ID:          filter.ID,
				ZoneID:      zone.ID,
				Description: filter.Description,
				Expression:  filter.Expression,
				Paused:      strconv.FormatBool(filter.Paused),
				Ref:         filter.Ref,
			},
			Meta:    make(map[string]string),
			Tainted: false,
		},
		DependsOn: []string{},
		Deposed:   []string{},
		Provider:  "provider.cloudflare",
		Type:      "cloudflare_filter",
	}

	return r
}
