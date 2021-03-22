package cmd

import (
	"github.com/spf13/cobra"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
)

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().StringVarP(&resourceType, "resource-type", "r", "", "Which resource you wish to generate")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Pull resources from the Cloudflare API and generate the respective Terraform resources",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("attempting to generating %q resources", *&resourceType)
		tmpDir, err := ioutil.TempDir("", "tfinstall")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		execPath, err := tfinstall.Find(context.Background(), tfinstall.LatestVersion(tmpDir, false))
		if err != nil {
			log.Fatal(err)
		}

		// Setup and configure Terraform to operate in the temporary directory where
		// the provider is already configured. Eventually, this will be '.'.
		workingDir := "/tmp"
		tf, err := tfexec.NewTerraform(workingDir, execPath)
		if err != nil {
			log.Fatal(err)
		}

		err = tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.LockTimeout("60s"))
		if err != nil {
			log.Fatal(err)
		}

		ps, err := tf.ProvidersSchema(context.Background())
		s := ps.Schemas["registry.terraform.io/cloudflare/cloudflare"]
		if s == nil {
			log.Fatal("failed to detect provider installation")
		}

		r := s.ResourceSchemas[*&resourceType]
	},
}
