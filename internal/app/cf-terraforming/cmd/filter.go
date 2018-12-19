package cmd

import (
	"fmt"
	"log"
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(filterCmd)
}

var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Import Filter data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("Importing Filter data")

		for _, zone := range zones {
			log.Printf("[DEBUG] Processing zone: ID %s, Name %s", zone.ID, zone.Name)

			filters, err := api.Filters(zone.ID, cloudflare.PaginationOptions{
				Page:    1,
				PerPage: 1000,
			})

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for _, r := range filters {
				log.Printf("[DEBUG] Filter ID %s, Expression %s, Description %s\n", r.ID, r.Expression, r.Description)
				// TODO: Process
			}
		}
	},
}
