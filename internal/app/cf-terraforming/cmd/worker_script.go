package cmd

import (
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go"

	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const workerScriptTemplate = `
resource "cloudflare_worker_script" "{{replace .ScriptName "-" "_"}}" {
	name    = "{{.ScriptName}}"
	content = <<-EOF
{{ trim .Script.Script }}
EOF
}
`

func init() {
	rootCmd.AddCommand(workerScriptCmd)
}

var workerScriptCmd = &cobra.Command{
	Use:   "worker_script",
	Short: "Import a worker script into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Importing worker script data")

		workerScripts, err := api.ListWorkerScripts()

		if err != nil {
			log.Error(err)
			return
		}
		// Loop through every script and fetch its content for rendering into tfstate format
		for _, script := range workerScripts.WorkerList {

			log.WithFields(logrus.Fields{
				"ID": script.ID,
			}).Debug("Fetching script")

			workerScriptResponse, downloadErr := api.DownloadWorker(&cloudflare.WorkerRequestParams{ScriptName: script.ID})

			if downloadErr != nil {
				if strings.Contains(downloadErr.Error(), "HTTP status 404") {
					log.WithFields(logrus.Fields{
						"Script ID": script.ID,
					}).Error("Error fetching script - does this zone have the workers entitlement?")
					continue
				}
				log.Error(downloadErr)
			}

			if workerScriptResponse.Success == true {

				log.WithFields(logrus.Fields{
					"Body": workerScriptResponse.WorkerScript,
				}).Debug("Worker script in multi-script mode")

				if tfstate {
					// TODO: Implement state dump
				} else {
					workerScriptParse(script.ID, workerScriptResponse.WorkerScript)
				}
			}
		}
	},
}

func workerScriptParse(scriptName string, script cloudflare.WorkerScript) {
	tmpl := template.Must(template.New("script").Funcs(templateFuncMap).Parse(workerScriptTemplate))
	err := tmpl.Execute(os.Stdout,
		struct {
			ScriptName string
			Script     cloudflare.WorkerScript
		}{
			ScriptName: scriptName,
			Script:     script,
		})
	if err != nil {
		log.Error(err)
	}
}
