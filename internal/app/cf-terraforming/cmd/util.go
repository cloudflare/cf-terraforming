package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

// testDataFile slurps a local test case into memory and returns it while
// encapsulating the logic for finding it.
func testDataFile(filename string) string {
	filename = strings.TrimSuffix(filename, "/")

	dirname, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dir, err := os.Open(filepath.Join(dirname, "../../../../testdata/terraform"))
	if err != nil {
		panic(err)
	}

	fullpath := dir.Name() + "/" + filename
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		panic(fmt.Errorf("terraform testdata file does not exist at %s", fullpath))
	}

	data, _ := ioutil.ReadFile(fullpath)

	return string(data)
}

func sharedPreRun(cmd *cobra.Command, args []string) {
	accountID = viper.GetString("account")

	if apiToken = viper.GetString("token"); apiToken == "" {
		if apiEmail = viper.GetString("email"); apiEmail == "" {
			log.Error("'email' must be set.")
		}

		if apiKey = viper.GetString("key"); apiKey == "" {
			log.Error("either -t/--token or -k/--key must be set.")
		}

		log.WithFields(logrus.Fields{
			"email":      apiEmail,
			"zone_id":    zoneID,
			"account_id": accountID,
		}).Debug("initializing cloudflare-go")

	} else {
		log.WithFields(logrus.Fields{
			"zone_id":    zoneID,
			"account_Id": accountID,
		}).Debug("initializing cloudflare-go with API Token")
	}

	var options []cloudflare.Option

	if accountID != "" {
		log.WithFields(logrus.Fields{
			"account_id": accountID,
		}).Debug("configuring Cloudflare API with account")

		// Organization ID was passed, use it to configure the API
		options = append(options, cloudflare.UsingAccount(accountID))
	}

	var err error

	// Don't initialise a client in CI as this messes with VCR and the ability to
	// mock out the HTTP interactions.
	if os.Getenv("CI") != "true" {
		var useToken = apiToken != ""

		if useToken {
			api, err = cloudflare.NewWithAPIToken(apiToken, options...)
		} else {
			api, err = cloudflare.New(apiKey, apiEmail, options...)
		}

		if err != nil {
			log.Fatal(err)
		}
	}
}

// sanitiseTerraformResourceName ensures that a Terraform resource name matches the
// restrictions imposed by core.
func sanitiseTerraformResourceName(s string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	return re.ReplaceAllString(s, "_")
}
