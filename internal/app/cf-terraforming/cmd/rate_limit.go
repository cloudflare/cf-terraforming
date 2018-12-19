package cmd

import (
	"fmt"
	"log"
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/spf13/cobra"
)

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
					// TODO: Process
				}
			}

		}
	},
}
