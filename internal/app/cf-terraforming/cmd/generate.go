package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().StringVarP(&resourceType, "resource-type", "r", "", "Which resource you wish to generate")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Pull resources from the Cloudflare API and generate the respective Terraform resources",
	Run: func(cmd *cobra.Command, args []string) {
		// you know, generate it
	},
}
