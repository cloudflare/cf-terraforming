package cmd

import (
	"fmt"
	"os"
	"strings"

	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const zoneLockdownTemplate = `
resource "cloudflare_zone_lockdown" "{{replace .Zone.Name "." "_"}}_{{.Lockdown.ID}}" {
    zone_id = "{{.Zone.ID}}"
    description = "{{.Lockdown.Description}}"
    urls = [
{{range .Lockdown.URLs}}
        "{{.}}",
{{end}}
    ]
    configurations = [
{{range .Lockdown.Configurations}}
        {
            target = "{{.Target}}"
            value = "{{.Value}}"
        },
{{end}}
    ]
}
`

func init() {
	rootCmd.AddCommand(zoneLockdownCmd)
}

var zoneLockdownCmd = &cobra.Command{
	Use:   "zone_lockdown",
	Short: "Import Zone Lockdown data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Importing Zone Lockdown data")

		for _, zone := range zones {

			log.WithFields(logrus.Fields{
				"ID":   zone.ID,
				"Name": zone.Name,
			}).Debug("Processing zone")

			totalPages := 999

			for page := 1; page <= totalPages; page++ {
				lockdowns, err := api.ListZoneLockdowns(zone.ID, page)

				if err != nil {
					log.Debug(err)
					return
				}

				totalPages = lockdowns.TotalPages

				for _, r := range lockdowns.Result {

					log.WithFields(logrus.Fields{
						"ID":  r.ID,
						"URL": r.URLs,
					}).Debug("Processing lockdown")

					if tfstate {
						r := zoneLockdownResourceStateBuild(zone, r)
						resourcesMap["cloudflare_zone_lockdown."+strings.ReplaceAll(zone.Name, ".", "_")+"_"+r.Primary.Id] = r
					} else {
						zoneLockdownParse(zone, r)
					}
				}
			}
		}
	},
}

func zoneLockdownParse(zone cloudflare.Zone, lockdown cloudflare.ZoneLockdown) {
	tmpl := template.Must(template.New("zone_lockdown").Funcs(templateFuncMap).Parse(zoneLockdownTemplate))
	err := tmpl.Execute(os.Stdout,
		struct {
			Zone     cloudflare.Zone
			Lockdown cloudflare.ZoneLockdown
		}{
			Zone:     zone,
			Lockdown: lockdown,
		})
	if err != nil {
		log.Error(err)
	}
}

func zoneLockdownResourceStateBuild(zone cloudflare.Zone, record cloudflare.ZoneLockdown) Resource {
	r := Resource{
		Primary: Primary{
			Id: record.ID,
			Attributes: map[string]string{
				"configurations.#": fmt.Sprint(len(record.Configurations)),
				"description":      record.Description,
				"id":               record.ID,
				"paused":           fmt.Sprint(record.Paused),
				"urls.#":           fmt.Sprint(len(record.URLs)),
				"zone_id":          zone.ID,
			},
			Meta:    make(map[string]string),
			Tainted: false,
		},
		DependsOn: []string{},
		Deposed:   []string{},
		Provider:  "provider.cloudflare",
		Type:      "cloudflare_zone_lockdown",
	}

	attributes := r.Primary.Attributes.(map[string]string)

	for _, configuration := range record.Configurations {
		hash := hashMap(map[string]string{
			"target": configuration.Target,
			"value":  configuration.Value,
		})
		attributes[fmt.Sprintf("configurations.%d.target", hash)] = configuration.Target
		attributes[fmt.Sprintf("configurations.%d.value", hash)] = configuration.Value
	}

	for _, url := range record.URLs {
		hash := hashcode.String(fmt.Sprintf("%s;", url))
		attributes[fmt.Sprintf("urls.%d", hash)] = url
	}

	r.Primary.Attributes = attributes

	return r
}
