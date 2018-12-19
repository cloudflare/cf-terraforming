package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(wafRuleCmd)
}

var wafRuleCmd = &cobra.Command{
	Use:   "waf_rule",
	Short: "Import WAF Rule data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("Importing WAF Rule data")

		for _, zone := range zones {
			log.Printf("[DEBUG] Processing zone: ID %s, Name %s", zone.ID, zone.Name)

			log.Printf("[DEBUG] Get WAF Packages")

			wafPackages, err := api.ListWAFPackages(zone.ID)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for _, wafPackage := range wafPackages {
				log.Printf("[DEBUG] WAF Package ID %s, Name %s\n", wafPackage.ID, wafPackage.Name)
				log.Printf("[DEBUG] Get WAF Rules in a Package")

				wafRules, err := api.ListWAFRules(zone.ID, wafPackage.ID)

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				for _, rule := range wafRules {
					log.Printf("[DEBUG] Processing WAF Rule: ID %s, Description %s", rule.ID, rule.Description)
					// TODO: Process
				}
			}
		}
	},
}
