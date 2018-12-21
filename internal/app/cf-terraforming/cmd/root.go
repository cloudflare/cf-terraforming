package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"

	cloudflare "github.com/cloudflare/cloudflare-go"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile, zoneName, apiEmail, apiKey, accountID string
var api *cloudflare.API
var zones []cloudflare.Zone

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cf-terraforming",
	Short: "Boostrapping Terraform from existing Cloudflare account",
	Long: `cf-terraforming is an application that allows Cloudflare users
to be able to adopt Terraform by giving them a feasible way to get
all of their existing Cloudflare configuration into Terraform.`,
	PersistentPreRun: persistentPreRun,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cf-terraforming.yaml)")

	// Zone selection
	rootCmd.PersistentFlags().StringVarP(&zoneName, "zone", "z", "", "Limit the export to a single zone (name or ID)")

	// Account
	rootCmd.PersistentFlags().StringVarP(&accountID, "account", "a", "", "Use specific account ID for import")

	// API credentials
	rootCmd.PersistentFlags().StringVar(&apiEmail, "email", "", "API Email address associated with your account")
	rootCmd.PersistentFlags().StringVar(&apiKey, "key", "", "API Key generated on the 'My Profile' page. See: https://dash.cloudflare.com/?account=profile")

	viper.BindPFlag("email", rootCmd.PersistentFlags().Lookup("email"))
	viper.BindPFlag("key", rootCmd.PersistentFlags().Lookup("key"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cf-terraforming" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cf-terraforming")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix("cf_terraforming")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func persistentPreRun(cmd *cobra.Command, args []string) {
	if apiEmail = viper.GetString("email"); apiEmail == "" {
		fmt.Println("'email' must be set.")
		os.Exit(1)
	}

	if apiKey = viper.GetString("key"); apiKey == "" {
		fmt.Println("'key' must be set.")
		os.Exit(1)
	}

	log.Printf("[DEBUG] API Email = %s\n", apiEmail)

	log.Print("[DEBUG] Initializing cloudflare-go")
	var err error
	api, err = cloudflare.New(apiKey, apiEmail)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[DEBUG] Selecting zones for import\n")

	if regexp.MustCompile("^[a-z0-9]{32}$").MatchString(zoneName) {
		zone, err := api.ZoneDetails(zoneName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		zones = []cloudflare.Zone{zone}
	} else if zoneName != "" {
		zones, err = api.ListZones(zoneName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		zones, err = api.ListZones()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	log.Printf("[INFO] Zones selected:\n")
	for _, i := range zones {
		log.Printf("[INFO] - ID: %s, Name: %s\n", i.ID, i.Name)
	}
}
