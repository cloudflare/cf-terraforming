package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(zoneLockdownCmd)
}

var zoneLockdownCmd = &cobra.Command{
	Use:   "zone_lockdown",
	Short: "Import Zone Lockdown data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("Importing Zone Lockdown data")

		for _, zone := range zones {
			log.Printf("[DEBUG] Processing zone: ID %s, Name %s", zone.ID, zone.Name)

			totalPages := 999

			for page := 1; page <= totalPages; page++ {
				lockdowns, err := api.ListZoneLockdowns(zone.ID, page)

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				totalPages = lockdowns.TotalPages

				for _, r := range lockdowns.Result {
					log.Printf("[DEBUG] Lockdown ID %s, URL %s\n", r.ID, r.URLs)
					// TODO: Process
				}
			}

		}
	},
}
