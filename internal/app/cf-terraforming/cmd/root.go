package cmd

import (
	"os"
	"regexp"
	"strings"

	"encoding/json"

	cloudflare "github.com/cloudflare/cloudflare-go"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var log = logrus.New()
var cfgFile, zoneName, apiEmail, apiKey, accountID, orgID, logLevel string
var verbose, tfstate bool
var api *cloudflare.API
var zones []cloudflare.Zone
var resourcesMap = map[string]interface{}{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cf-terraforming",
	Short: "Boostrapping Terraform from existing Cloudflare account",
	Long: `cf-terraforming is an application that allows Cloudflare users
to be able to adopt Terraform by giving them a feasible way to get
all of their existing Cloudflare configuration into Terraform.`,
	PersistentPreRun:  persistentPreRun,
	PersistentPostRun: persistentPostRun,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Debug(err)
		return
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.cf-terraforming.yaml)")

	// Zone selection
	rootCmd.PersistentFlags().StringVarP(&zoneName, "zone", "z", "", "Limit the export to a single zone (name or ID)")

	// Account
	rootCmd.PersistentFlags().StringVarP(&accountID, "account", "a", "", "Use specific account ID for import")

	// API credentials
	rootCmd.PersistentFlags().StringVarP(&apiEmail, "email", "e", "", "API Email address associated with your account")
	rootCmd.PersistentFlags().StringVarP(&apiKey, "key", "k", "", "API Key generated on the 'My Profile' page. See: https://dash.cloudflare.com/?account=profile")

	// [Optional] Organization ID
	rootCmd.PersistentFlags().StringVarP(&orgID, "organization", "o", "", "Use specific organization ID for import")

	// Debug logging mode
	rootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "", "Specify logging level: (trace, debug, info, warn, error, fatal, panic)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Specify verbose output (same as setting log level to debug)")

	// tfstate export
	rootCmd.PersistentFlags().BoolVarP(&tfstate, "tfstate", "s", false, "Export tfstate for the given resource instead of HCL Terraform config (default)")

	viper.BindPFlag("email", rootCmd.PersistentFlags().Lookup("email"))
	viper.BindPFlag("key", rootCmd.PersistentFlags().Lookup("key"))
	viper.BindPFlag("organization", rootCmd.PersistentFlags().Lookup("organization"))
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
		log.Debug("Using config file:", viper.ConfigFileUsed())
	}

	var cfgLogLevel = logrus.InfoLevel

	// A user may also pass the verbose flag in order to support this convention
	if verbose {
		cfgLogLevel = logrus.DebugLevel
	}

	switch strings.ToLower(logLevel) {
	case "trace":
		cfgLogLevel = logrus.TraceLevel
	case "debug":
		cfgLogLevel = logrus.DebugLevel
	case "info":
		break
	case "warn":
		cfgLogLevel = logrus.WarnLevel
	case "error":
		cfgLogLevel = logrus.ErrorLevel
	case "fatal":
		cfgLogLevel = logrus.FatalLevel
	case "panic":
		cfgLogLevel = logrus.PanicLevel
	}

	log.SetLevel(cfgLogLevel)
}

func persistentPreRun(cmd *cobra.Command, args []string) {
	if apiEmail = viper.GetString("email"); apiEmail == "" {
		log.Error("'email' must be set.")
	}

	if apiKey = viper.GetString("key"); apiKey == "" {
		log.Error("'key' must be set.")
	}

	log.WithFields(logrus.Fields{
		"API email":       apiEmail,
		"Zone name":       zoneName,
		"Account ID":      accountID,
		"Organization ID": orgID,
	}).Debug("Initializing cloudflare-go")

	var options cloudflare.Option

	if orgID = viper.GetString("organization"); orgID != "" {
		log.WithFields(logrus.Fields{
			"ID": orgID,
		}).Debug("Configuring Cloudflare API with organization")

		// Organization ID was passed, use it to configure the API
		options = cloudflare.UsingOrganization(orgID)
	}

	var err error
	if options != nil {
		api, err = cloudflare.New(apiKey, apiEmail, options)
	} else {
		api, err = cloudflare.New(apiKey, apiEmail)
	}
	if err != nil {
		log.Error(err)
	}

	log.Debug("Selecting zones for import")

	if regexp.MustCompile("^[a-z0-9]{32}$").MatchString(zoneName) {
		zone, err := api.ZoneDetails(zoneName)

		if err != nil {
			log.Error(err)
		}

		zones = []cloudflare.Zone{zone}
	} else if zoneName != "" {
		zones, err = api.ListZones(zoneName)

		if err != nil {
			log.Error(err)
		}
	} else {
		zones, err = api.ListZones()

		if err != nil {
			log.Error(err)
		}
	}

	log.Debug("Zones selected:\n")

	for _, i := range zones {

		log.WithFields(logrus.Fields{
			"ID":   i.ID,
			"Name": i.Name,
		}).Debug("Zone")
	}
}

// This function runs following every root command
// When the tfstate flag is passed, output the
// full Terraform state file by rendering the resources map
// that was built up by the resource-specific commands
func persistentPostRun(cmd *cobra.Command, args []string) {

	if tfstate {

		m := []Module{{
			Path:      []string{"root"},
			DependsOn: []string{},
			Outputs:   make(map[string]string),
			Resource:  resourcesMap,
		}}

		s := TFStateScaffold{
			Version: 1,
			Serial:  0,
			Modules: m,
		}

		sr := StateResponse{
			TFStateScaffold: s,
		}

		j, _ := json.Marshal(sr)

		os.Stdout.Write(j)
	}
}
