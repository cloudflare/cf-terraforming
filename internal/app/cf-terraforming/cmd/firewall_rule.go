package cmd

import (
	"fmt"
	"log"
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(firewallRuleCmd)
}

var firewallRuleCmd = &cobra.Command{
	Use:   "firewall_rule",
	Short: "Import Firewall Rule data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("Importing Firewall Rule data")

		for _, zone := range zones {
			log.Printf("[DEBUG] Processing zone: ID %s, Name %s", zone.ID, zone.Name)

			firewallRules, err := api.FirewallRules(zone.ID, cloudflare.PaginationOptions{
				Page:    1,
				PerPage: 1000,
			})

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for _, r := range firewallRules {
				log.Printf("[DEBUG] Firewall Rule ID %s, Description %s\n", r.ID, r.Description)
				// TODO: Process
			}
		}
	},
}
