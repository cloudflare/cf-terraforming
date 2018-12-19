package cmd

import (
	"fmt"
	"log"
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(recordCmd)
}

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Import Record data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("Importing DNS Record data")

		for _, zone := range zones {
			log.Printf("[DEBUG] Processing zone: ID %s, Name %s", zone.ID, zone.Name)

			// Fetch all records for a zone
			recs, err := api.DNSRecords(zone.ID, cloudflare.DNSRecord{})

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			for _, r := range recs {
				fmt.Printf("Record ID %s, Name %s, Type %s: %s\n", r.ID, r.Name, r.Type, r.Content)
			}
		}
	},
}
