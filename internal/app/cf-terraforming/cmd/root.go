package cmd

import (
	cloudflare "github.com/cloudflare/cloudflare-go"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var log = logrus.New()
var cfgFile, zoneID, hostname, apiEmail, apiKey, apiToken, accountID, terraformInstallPath string
var verbose bool
var api *cloudflare.API
var terraformImportCmdPrefix = "terraform import"
var terraformResourceNamePrefix = "terraform_managed_resource"

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "cf-terraforming",
	Short: "Bootstrapping Terraform from existing Cloudflare account",
	Long: `cf-terraforming is an application that allows Cloudflare users
to be able to adopt Terraform by giving them a feasible way to get
all of their existing Cloudflare configuration into Terraform.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		return
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	home, err := homedir.Dir()
	if err != nil {
		log.Debug(err)
		return
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", home+"/.cf-terraforming.yaml", "Path to config file")

	// Zone selection
	rootCmd.PersistentFlags().StringVarP(&zoneID, "zone", "z", "", "Limit the export to a single zone ID")
	if err = viper.BindPFlag("zone", rootCmd.PersistentFlags().Lookup("zone")); err != nil {
		log.Fatal(err)
	}
	if err = viper.BindEnv("zone", "CLOUDFLARE_ZONE_ID"); err != nil {
		log.Fatal(err)
	}

	// Account
	rootCmd.PersistentFlags().StringVarP(&accountID, "account", "a", "", "Use specific account ID for commands")
	if err = viper.BindPFlag("account", rootCmd.PersistentFlags().Lookup("account")); err != nil {
		log.Fatal(err)
	}
	if err = viper.BindEnv("account", "CLOUDFLARE_ACCOUNT_ID"); err != nil {
		log.Fatal(err)
	}

	// API credentials
	rootCmd.PersistentFlags().StringVarP(&apiEmail, "email", "e", "", "API Email address associated with your account")
	if err = viper.BindPFlag("email", rootCmd.PersistentFlags().Lookup("email")); err != nil {
		log.Fatal(err)
	}
	if err = viper.BindEnv("email", "CLOUDFLARE_EMAIL"); err != nil {
		log.Fatal(err)
	}

	rootCmd.PersistentFlags().StringVarP(&apiKey, "key", "k", "", "API Key generated on the 'My Profile' page. See: https://dash.cloudflare.com/profile/api-tokens")
	if err = viper.BindPFlag("key", rootCmd.PersistentFlags().Lookup("key")); err != nil {
		log.Fatal(err)
	}
	if err = viper.BindEnv("key", "CLOUDFLARE_API_KEY"); err != nil {
		log.Fatal(err)
	}

	rootCmd.PersistentFlags().StringVarP(&apiToken, "token", "t", "", "API Token")
	if err = viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token")); err != nil {
		log.Fatal(err)
	}
	if err = viper.BindEnv("token", "CLOUDFLARE_API_TOKEN"); err != nil {
		log.Fatal(err)
	}

	rootCmd.PersistentFlags().StringVarP(&hostname, "hostname", "", "", "Hostname to use to query the API")
	if err = viper.BindPFlag("hostname", rootCmd.PersistentFlags().Lookup("hostname")); err != nil {
		log.Fatal(err)
	}
	if err = viper.BindEnv("hostname", "CLOUDFLARE_API_HOSTNAME"); err != nil {
		log.Fatal(err)
	}

	// Debug logging mode
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Specify verbose output (same as setting log level to debug)")

	rootCmd.PersistentFlags().StringVar(&resourceType, "resource-type", "", "Which resource you wish to generate")

	rootCmd.PersistentFlags().StringVar(&terraformInstallPath, "terraform-install-path", ".", "Path to the Terraform installation")

	if err = viper.BindPFlag("terraform-install-path", rootCmd.PersistentFlags().Lookup("terraform-install-path")); err != nil {
		log.Fatal(err)
	}
	if err = viper.BindEnv("terraform-install-path", "CLOUDFLARE_TERRAFORM_INSTALL_PATH"); err != nil {
		log.Fatal(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Debug(err)
			return
		}

		// Search config in home directory with name ".cf-terraforming" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cf-terraforming")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix("cf_terraforming")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debug("using config file:", viper.ConfigFileUsed())
	}

	// var cfgLogLevel = logrus.InfoLevel

	// if verbose {
	// 	cfgLogLevel = logrus.DebugLevel
	// }
	var cfgLogLevel = logrus.InfoLevel
	cfgLogLevel = logrus.DebugLevel

	log.SetLevel(cfgLogLevel)
}
