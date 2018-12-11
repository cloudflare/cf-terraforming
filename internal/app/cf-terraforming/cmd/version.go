package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionString = "dev"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of cf-terraforming",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cf-terraforming version", versionString)
	},
}
