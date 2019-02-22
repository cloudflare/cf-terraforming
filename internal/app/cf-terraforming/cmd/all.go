package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(allCmd)
}

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Import all Cloudflare resources into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Importing ALL Cloudflare resource data")

		accessApplicationCmd.Run(cmd, args)

		accessPolicyCmd.Run(cmd, args)

		accessRuleCmd.Run(cmd, args)

		accountMemberCmd.Run(cmd, args)

		customPagesCmd.Run(cmd, args)

		filterCmd.Run(cmd, args)

		firewallRuleCmd.Run(cmd, args)

		loadBalancerCmd.Run(cmd, args)

		loadBalancerMonitorCmd.Run(cmd, args)

		loadBalancerPoolCmd.Run(cmd, args)

		pageRuleCmd.Run(cmd, args)

		rateLimitCmd.Run(cmd, args)

		recordCmd.Run(cmd, args)

		spectrumApplicationCmd.Run(cmd, args)

		wafRuleCmd.Run(cmd, args)

		workerRouteCmd.Run(cmd, args)

		workerScriptCmd.Run(cmd, args)

		zoneCmd.Run(cmd, args)

		zoneLockdownCmd.Run(cmd, args)

		zoneSettingsOverrideCmd.Run(cmd, args)
	},
}
