package cmd

import (
	"fmt"
	"os"
	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const recordTemplate = `
resource "cloudflare_record" "{{.Record.Type}}_{{replace .Record.Name "." "_"}}_{{.Record.ID}}" {
    domain = "{{.Zone.Name}}"
{{ if .Zone.Paused}}
    paused = "true"
{{end}}
    name = "{{.Record.Name}}"
    type = "{{.Record.Type}}"
    ttl = "{{.Record.TTL}}"
    proxied = "{{.Record.Proxied}}"
{{ if or (eq .Record.Type "MX") (eq .Record.Type "URI") }}
    priority = "{{.Record.Priority}}"
{{end}}
{{ if .IsValueTypeField }}
    value = "{{.Record.Content}}"
{{end}}
{{ if .IsDataTypeField }}
    data {
{{range $k, $v := .Record.Data}}
        {{ $k }} = "{{ $v }}"
{{end}}
    }
{{end}}
}

`

var dnsTypeValueFields = []string{
	"A", "AAAA", "CNAME", "NS", "MX", "TXT", "SPF",
}

var dnsTypeDataFields = []string{
	"LOC", "SRV", "CERT", "DNSKEY", "DS", "NAPTR", "SMIMEA", "SSHFP", "TLSA", "URI",
}

func init() {
	rootCmd.AddCommand(recordCmd)
}

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Import Record data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Importing DNS Record data")

		for _, zone := range zones {

			log.WithFields(logrus.Fields{
				"ID":   zone.ID,
				"Name": zone.Name,
			}).Debug("Processing zone")

			// Fetch all records for a zone
			recs, err := api.DNSRecords(zone.ID, cloudflare.DNSRecord{})

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			for _, r := range recs {

				log.WithFields(logrus.Fields{
					"ID":      r.ID,
					"Name":    r.Name,
					"Type":    r.Type,
					"Content": r.Content,
				}).Debug("Processing record")

				recordParse(zone, r)
			}
		}
	},
}

func recordParse(zone cloudflare.Zone, record cloudflare.DNSRecord) {
	tmpl := template.Must(template.New("record").Funcs(templateFuncMap).Parse(recordTemplate))
	if err := tmpl.Execute(os.Stdout,
		struct {
			Zone             cloudflare.Zone
			Record           cloudflare.DNSRecord
			IsValueTypeField bool
			IsDataTypeField  bool
		}{
			Zone:             zone,
			Record:           record,
			IsValueTypeField: contains(dnsTypeValueFields, record.Type),
			IsDataTypeField:  contains(dnsTypeDataFields, record.Type),
		}); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
