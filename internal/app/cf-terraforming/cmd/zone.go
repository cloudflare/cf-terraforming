package cmd

import (
	"fmt"
	"log"
	"os"

	// cloudflare "github.com/cloudflare/cloudflare-go"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spf13/cobra"
	terraformProviderCloudflare "github.com/terraform-providers/terraform-provider-cloudflare/cloudflare"
)

func init() {
	rootCmd.AddCommand(zoneCmd)
}

var zoneCmd = &cobra.Command{
	Use:   "zone",
	Short: "Import zone data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("Importing zones' data")
		log.Printf("%#v", zones)

		provider := terraformProviderCloudflare.Provider()
		resource := provider.(*schema.Provider).ResourcesMap["cloudflare_zone"]

		log.Printf("resourceCloudflareZone: %#v", provider)
		log.Printf("%#v", resource)

		for _, zone := range zones {
			zoneDetails, err := api.ZoneDetails(zone.ID)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			log.Printf("[DEBUG] Processing zone: ID %s, Name %s", zoneDetails.ID, zoneDetails.Name)

			resourceData := &schema.ResourceData{}
			resource.Read(resourceData, api)
			log.Printf("resourceData %#v", resourceData)

			// log.Printf("access policy %#v", cloudflare.AccessPolicy)
		}
	},
}
