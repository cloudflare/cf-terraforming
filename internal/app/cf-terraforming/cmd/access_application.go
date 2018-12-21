package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(accessApplicationCmd)
}

var accessApplicationCmd = &cobra.Command{
	Use:   "access_application",
	Short: "Import Access Application data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("Importing Access Application data")

		for _, zone := range zones {
			log.Printf("[DEBUG] Processing zone: ID %s, Name %s", zone.ID, zone.Name)

			accessApplications, _, err := api.AccessApplications(zone.ID, cloudflare.PaginationOptions{
				Page:    1,
				PerPage: 1000,
			})

			if err != nil {
				if strings.Contains(err.Error(), "HTTP status 403") {
					log.Printf("[INFO] insufficient permissions accessing Zone ID %s\n", zone.ID)
					continue
				}

				fmt.Println(err)
				os.Exit(1)
			}

			for _, r := range accessApplications {
				log.Printf("[DEBUG] Access Application ID %s, Name %s, Domain %s\n", r.ID, r.Name, r.Domain)
				// TODO: Process
			}

		}
	},
}
