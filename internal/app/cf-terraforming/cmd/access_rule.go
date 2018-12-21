package cmd

import (
	"fmt"
	"log"
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(accessRuleCmd)
}

var accessRuleCmd = &cobra.Command{
	Use:   "access_rule",
	Short: "Import Access Rule data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("Importing Access Rule data")

		for _, zone := range zones {
			log.Printf("[DEBUG] Processing zone: ID %s, Name %s", zone.ID, zone.Name)

			totalPages := 999

			for page := 1; page <= totalPages; page++ {
				accessRules, err := api.ListZoneAccessRules(zone.ID, cloudflare.AccessRule{}, page)

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				totalPages = accessRules.TotalPages

				for _, r := range accessRules.Result {
					log.Printf("[DEBUG] Filter ID %s, Notes %s, Configuration %s, Scope %s\n", r.ID, r.Notes, r.Configuration, r.Scope)
					// TODO: Process
				}
			}
		}
	},
}
