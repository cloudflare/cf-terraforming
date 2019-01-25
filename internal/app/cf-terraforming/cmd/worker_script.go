package cmd

import (
	"fmt"
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go"

	"text/template"

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
		// Enterprise multi-script mode:
		// If the organization ID is set at the API level,
		// this is an enterprise request, which can use special endpoints such as enumerate workers
		if api.OrganizationID != "" {
			workerScripts, err := api.ListWorkerScripts()

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			// Loop through every script and fetch its content for rendering into tfstate format
			for _, script := range workerScripts.WorkerList {

				workerScriptResponse, downloadErr := api.DownloadWorker(&cloudflare.WorkerRequestParams{ScriptName: script.ID})

				if downloadErr != nil {
					fmt.Println(downloadErr)
					os.Exit(1)
				}

				if workerScriptResponse.Success == true {
					workerScriptParse(script.ID, "", workerScriptResponse.WorkerScript, true)
				}
			}

		} else {
			for _, zone := range zones {
				// Otherwise, we're dealing with single script mode for all other zone types
				workerScriptResponse, singleScriptErr := api.DownloadWorker(&cloudflare.WorkerRequestParams{ZoneID: zone.ID})

				if singleScriptErr != nil {
					fmt.Println(singleScriptErr)
					os.Exit(1)
				}
				// It's possible for the script ID to be unset in some cases,
				// so set a default value
				var scriptID = "my_script"
				if workerScriptResponse.WorkerScript.WorkerMetaData.ID != "" {
					scriptID = workerScriptResponse.WorkerScript.WorkerMetaData.ID
				}

				if workerScriptResponse.Success == true {
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
