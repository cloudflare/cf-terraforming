package cmd

import (
	"os"
	"strconv"
	"strings"

	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const zoneTemplate = `
resource "cloudflare_zone" "{{replace .Zone.Name "." "_"}}" {
    zone = "{{.Zone.Name}}"
{{ if .Zone.Paused}}    paused = "true"{{end}}
    plan = "{{.ZonePlan}}"
}

`

type ZoneAttributes struct {
	ID                   string `json:"id"`
	PhishingDetected     string `json:"meta.phishing_detected"`
	WildcardProxiable    string `json:"meta.wildcard_proxiable"`
	NameServersNum       string `json:"name_servers.#"`
	NameServers0         string `json:"name_servers.0"`
	NameServers1         string `json:"name_servers.1"`
	Paused               string `json:"paused"`
	Plan                 string `json:"plan"`
	Status               string `json:"status"`
	Type                 string `json:"type"`
	VanityNameServersNum string `json:"vanity_name_servers.#"`
	VanityNameServers0   string `json:"vanity_name_servers.0"`
	VanityNameServers1   string `json:"vanity_name_servers.1"`
	Zone                 string `json:"zone"`
}

// we enforce the use of the Cloudflare API 'legacy_id' field until the mapping of plan is fixed in cloudflare-go
const (
	planIDFree       = "free"
	planIDPro        = "pro"
	planIDBusiness   = "business"
	planIDEnterprise = "enterprise"
)

// we keep a private map and we will have a function to check and validate the descriptive name from the RatePlan API with the legacy_id
var idForName = map[string]string{
	"Free Website":       planIDFree,
	"Pro Website":        planIDPro,
	"Business Website":   planIDBusiness,
	"Enterprise Website": planIDEnterprise,
}

func init() {
	rootCmd.AddCommand(zoneCmd)
}

var zoneCmd = &cobra.Command{
	Use:   "zone",
	Short: "Import zone data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Importing zone data")

		for _, zone := range zones {
			zoneDetails, err := api.ZoneDetails(zone.ID)

			if err != nil {
				log.Debug(err)
				return
			}

			log.WithFields(logrus.Fields{
				"ID":   zoneDetails.ID,
				"Name": zoneDetails.Name,
			}).Debug("Processing zone")

			if tfstate {
				r := zoneResourceStateBuild(zone)
				resourcesMap["cloudflare_zone."+strings.ReplaceAll(zoneDetails.Name, ".", "_")] = r
			} else {
				zoneParse(zone)
			}
		}
	},
}

func zoneParse(zone cloudflare.Zone) {
	tmpl := template.Must(template.New("zone").Funcs(templateFuncMap).Parse(zoneTemplate))
	err := tmpl.Execute(os.Stdout,
		struct {
			Zone     cloudflare.Zone
			ZonePlan string
		}{
			Zone:     zone,
			ZonePlan: idForName[zone.Plan.Name],
		})
	if err != nil {
		log.Error(err)
	}
}

func zoneResourceStateBuild(zone cloudflare.Zone) Resource {
	r := Resource{
		Primary: Primary{
			Id: zone.ID,
			Attributes: ZoneAttributes{
				ID:                   zone.ID,
				PhishingDetected:     strconv.FormatBool(zone.Meta.PhishingDetected),
				WildcardProxiable:    strconv.FormatBool(zone.Meta.WildcardProxiable),
				NameServersNum:       strconv.Itoa(len(zone.NameServers)),
				NameServers0:         zone.NameServers[0],
				NameServers1:         zone.NameServers[1],
				Paused:               strconv.FormatBool(zone.Paused),
				Plan:                 idForName[zone.Plan.Name],
				Status:               zone.Status,
				VanityNameServersNum: strconv.Itoa(len(zone.VanityNS)),
				Type:                 zone.Type,
				Zone:                 zone.Name,
			},
			Meta:    make(map[string]string),
			Tainted: false,
		},
		DependsOn: []string{},
		Deposed:   []string{},
		Provider:  "provider.cloudflare",
		Type:      "cloudflare_zone",
	}

	zoneAttributes := r.Primary.Attributes.(ZoneAttributes)
	if len(zone.VanityNS) > 0 {
		zoneAttributes.VanityNameServers0 = zone.VanityNS[0]
	}
	if len(zone.VanityNS) > 1 {
		zoneAttributes.VanityNameServers1 = zone.VanityNS[0]
	}
	if len(zone.NameServers) > 0 {
		zoneAttributes.NameServers0 = zone.NameServers[0]
	}
	if len(zone.NameServers) > 1 {
		zoneAttributes.NameServers1 = zone.NameServers[1]
	}
	r.Primary.Attributes = zoneAttributes

	return r
}
