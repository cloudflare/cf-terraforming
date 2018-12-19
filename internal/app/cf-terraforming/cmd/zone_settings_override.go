package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(zoneSettingsOverrideCmd)
}

var zoneSettingsOverrideCmd = &cobra.Command{
	Use:   "zone_settings_override",
	Short: "Import Zone Settings Override data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("Importing Zone Settings data")

		for _, zone := range zones {
			log.Printf("[DEBUG] Processing zone: ID %s, Name %s", zone.ID, zone.Name)

			// Fetch all records for a zone
			settings, err := api.ZoneSettings(zone.ID)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			for _, s := range settings.Result {
				fmt.Printf("Setting ID %s, Value %s\n", s.ID, s.Value)
				// Process
			}
		}
	},
}
