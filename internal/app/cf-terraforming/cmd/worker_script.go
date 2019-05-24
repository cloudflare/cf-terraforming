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
{{if not .MultiScript }}
	zone = "{{.ZoneName}}"
{{else}}
	name  = "{{replace .ScriptName "-" "_"}}"
{{end}}
    content = "{{replace .Script.Script "\n" ""}}"
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
		// Enterprise multi-script mode:
		// If the organization ID is set at the API level,
		// this is an enterprise request, which can use special endpoints such as enumerate workers
		if api.OrganizationID != "" {
			workerScripts, err := api.ListWorkerScripts()

			if err != nil {
				log.Debug(err)
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
						}).Debug("Error fetching script - does this zone have the workers entitlement?")

						continue
					}
					log.Debug(downloadErr)
				}

				if workerScriptResponse.Success == true {

					log.WithFields(logrus.Fields{
						"Body": workerScriptResponse.WorkerScript,
					}).Debug("Worker script in multi-script mode")

					if tfstate {
						// TODO: Implement state dump
					} else {
						workerScriptParse(script.ID, "", workerScriptResponse.WorkerScript, true)
					}
				}
			}

		} else {
			for _, zone := range zones {
				// Otherwise, we're dealing with single script mode for all other zone types
				workerScriptResponse, singleScriptErr := api.DownloadWorker(&cloudflare.WorkerRequestParams{ZoneID: zone.ID})

				if singleScriptErr != nil {
					//Workers endpoints may return a 404 if the zone is not entitled to use workers
					//skip over this error to avoid polluting stdout / generated config files
					if strings.Contains(singleScriptErr.Error(), "HTTP status 404") {
						log.WithFields(logrus.Fields{
							"Zone ID": zone.ID,
						}).Debug("Error fetching script - does this zone have the workers entitlement?")

						continue
					}
					log.Debug(singleScriptErr)
				}
				// It's possible for the script ID to be unset in some cases,
				// so set a default value
				var scriptID = "my_script"
				if workerScriptResponse.WorkerScript.WorkerMetaData.ID != "" {
					scriptID = workerScriptResponse.WorkerScript.WorkerMetaData.ID
				}

				if workerScriptResponse.Success == true {

					log.WithFields(logrus.Fields{
						"Body": workerScriptResponse.WorkerScript,
					}).Debug("Worker script in non multi-script mode")

					workerScriptParse(scriptID, zone.Name, workerScriptResponse.WorkerScript, false)
				}
			}
		}
	},
}

func workerScriptParse(scriptName string, zoneName string, script cloudflare.WorkerScript, multiScript bool) {
	tmpl := template.Must(template.New("script").Funcs(templateFuncMap).Parse(workerScriptTemplate))
	tmpl.Execute(os.Stdout,
		struct {
			ScriptName  string
			ZoneName    string
			Script      cloudflare.WorkerScript
			MultiScript bool
		}{
			ScriptName:  scriptName,
			ZoneName:    zoneName,
			Script:      script,
			MultiScript: multiScript,
		})
}
