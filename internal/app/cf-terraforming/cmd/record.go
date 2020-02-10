package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const recordTemplate = `
resource "cloudflare_record" "{{.Record.Type}}_{{replace .Record.Name "." "_"}}_{{.Record.ID}}" {
    zone_id = "{{.Zone.ID}}"
{{ if .Zone.Paused}}
    paused = "true"
{{end}}
    name = "{{normalizeRecordName .Record.Name .Record.ZoneName}}"
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
    data = {
{{range $k, $v := .Record.Data}}
        {{ $k }} = "{{ $v }}"
{{end}}
    }
{{end}}
}

`

type RecordAttributes struct {
	ID                          string `json:"id"`
	CreatedOn                   string `json:"created_on"`
	DataNum                     string `json:"data.%"`
	DataAlgorithm               string `json:"data.algorithm,omitempty"`
	DataCertificate             string `json:"data.certificate,omitempty"`
	DataContent                 string `json:"data.content,omitempty"`
	DataDigest                  string `json:"data.digest,omitempty"`
	DataDigestType              string `json:"data.digest_type,omitempty"`
	DataFingerprint             string `json:"data.fingerprint,omitempty"`
	DataFlags                   string `json:"data.flags,omitempty"`
	DataOrder                   string `json:"data.order,omitempty"`
	DataKeyTag                  string `json:"data.key_tag,omitempty"`
	DataMatchingType            string `json:"data.matching_type,omitempty"`
	DataName                    string `json:"data.name,omitempty"`
	DataPort                    string `json:"data.port,omitempty"`
	DataPreference              string `json:"data.preference,omitempty"`
	DataPriority                string `json:"data.priority,omitempty"`
	DataProto                   string `json:"data.proto,omitempty"`
	DataProtocol                string `json:"data.protocol,omitempty"`
	DataPublicKey               string `json:"public_key,omitempty"`
	DataRegex                   string `json:"data.regex,omitempty"`
	DataReplacement             string `json:"data.replacement,omitempty"`
	DataSelector                string `json:"data.selector,omitempty"`
	DataService                 string `json:"data.service,omitempty"`
	DataTag                     string `json:"data.tag,omitempty"`
	DataTarget                  string `json:"data.target,omitempty"`
	DataType                    string `json:"data.type,omitempty"`
	DataUsage                   string `json:"data.usage,omitempty"`
	DataValue                   string `json:"data.value,omitempty"`
	DataWeight                  string `json:"data.weight,omitempty"`
	DataAltitude                string `json:"data.altitude,omitempty"`
	DataLatDegrees              string `json:"data.lat_degrees,omitempty"`
	DataLatDirection            string `json:"data.lat_direction,omitempty"`
	DataLatMinutes              string `json:"data.lat_minutes,omitempty"`
	DataLatSeconds              string `json:"data.lat_seconds,omitempty"`
	DataLongDegrees             string `json:"data.long_degrees,omitempty"`
	DataLongDirection           string `json:"data.long_direction,omitempty"`
	DataLongMinutes             string `json:"data.long_minutes,omitempty"`
	DataLongSeconds             string `json:"data.long_seconds,omitempty"`
	DataPrecisionHorz           string `json:"data.precision_horz,omitempty"`
	DataPrecisionVert           string `json:"data.precision_vert,omitempty"`
	DataSize                    string `json:"data.size,omitempty"`
	Domain                      string `json:"domain"`
	Hostname                    string `json:"hostname"`
	MetadataNum                 string `json:"metadata.%"`
	MetadataAutoAdded           string `json:"metadata.auto_added,omitempty"`
	MetadataManagedByApps       string `json:"metadata.managed_by_apps,omitempty"`
	MetadataManagedByArgoTunnel string `json:"metadata.managed_by_argo_tunnel,omitempty"`
	ModifiedOn                  string `json:"modified_on"`
	Name                        string `json:"name"`
	Priority                    string `json:"priority"`
	Proxiable                   string `json:"proxiable"`
	Proxied                     string `json:"proxied"`
	TTL                         string `json:"ttl"`
	Type                        string `json:"type"`
	Value                       string `json:"value"`
	ZoneID                      string `json:"zone_id"`
}

var dnsTypeValueFields = []string{
	"A", "AAAA", "CNAME", "NS", "MX", "TXT", "SPF",
}

var dnsTypeDataFields = []string{
	"LOC", "SRV", "CAA", "CERT", "DNSKEY", "DS", "NAPTR", "SMIMEA", "SSHFP", "TLSA", "URI",
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
				log.Debug(err)
				return
			}
			for _, r := range recs {

				log.WithFields(logrus.Fields{
					"ID":      r.ID,
					"Name":    r.Name,
					"Type":    r.Type,
					"Content": r.Content,
				}).Debug("Processing record")

				if tfstate {
					state := recordResourceStateBuild(zone, r)
					name := r.Type + "_" + strings.ReplaceAll(r.Name, ".", "_") + "_" + r.ID
					resourcesMap["cloudflare_record."+name] = state
				} else {
					recordParse(zone, r)
				}
			}
		}
	},
}

func recordParse(zone cloudflare.Zone, record cloudflare.DNSRecord) {
	tmpl := template.Must(template.New("record").Funcs(templateFuncMap).Parse(recordTemplate))
	tmpl.Execute(os.Stdout,
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
		})
}

func recordResourceStateBuild(zone cloudflare.Zone, record cloudflare.DNSRecord) Resource {
	var meta map[string]interface{}
	if record.Meta != nil {
		meta = record.Meta.(map[string]interface{})
	}

	var data map[string]interface{}
	if record.Data != nil {
		data = record.Data.(map[string]interface{})
	}

	log.Printf("%#v", data)

	r := Resource{
		Primary: Primary{
			Id: record.ID,
			Attributes: RecordAttributes{
				ID:                          record.ID,
				CreatedOn:                   record.CreatedOn.Format(time.RFC3339),
				DataNum:                     strconv.Itoa(len(data)),
				Domain:                      record.ZoneName,
				Hostname:                    record.Name,
				MetadataNum:                 strconv.Itoa(len(record.Meta.(map[string]interface{}))),
				MetadataAutoAdded:           strconv.FormatBool(meta["auto_added"].(bool)),
				MetadataManagedByApps:       strconv.FormatBool(meta["managed_by_apps"].(bool)),
				MetadataManagedByArgoTunnel: strconv.FormatBool(meta["managed_by_argo_tunnel"].(bool)),
				ModifiedOn:                  record.ModifiedOn.Format(time.RFC3339),
				Name:                        normalizeRecordName(record.Name, record.ZoneName),
				Priority:                    strconv.Itoa(record.Priority),
				Proxiable:                   strconv.FormatBool(record.Proxiable),
				Proxied:                     strconv.FormatBool(record.Proxied),
				TTL:                         strconv.Itoa(record.TTL),
				Type:                        record.Type,
				Value:                       record.Content,
				ZoneID:                      record.ZoneID,
			},
			Meta:    make(map[string]string),
			Tainted: false,
		},
		DependsOn: []string{},
		Deposed:   []string{},
		Provider:  "provider.cloudflare",
		Type:      "cloudflare_record",
	}

	recordAttributes := r.Primary.Attributes.(RecordAttributes)
	if v, ok := data["algorithm"]; ok {
		recordAttributes.DataAlgorithm = fmt.Sprint(v)
	}
	if v, ok := data["certificate"]; ok {
		recordAttributes.DataCertificate = v.(string)
	}
	if v, ok := data["content"]; ok {
		recordAttributes.DataContent = v.(string)
	}
	if v, ok := data["digest"]; ok {
		recordAttributes.DataDigest = v.(string)
	}
	if v, ok := data["digest_type"]; ok {
		recordAttributes.DataDigestType = fmt.Sprint(v)
	}
	if v, ok := data["fingerprint"]; ok {
		recordAttributes.DataFingerprint = v.(string)
	}
	if v, ok := data["flags"]; ok {
		recordAttributes.DataFlags = fmt.Sprint(v)
	}
	if v, ok := data["key_tag"]; ok {
		recordAttributes.DataKeyTag = fmt.Sprint(v)
	}
	if v, ok := data["matching_type"]; ok {
		recordAttributes.DataMatchingType = fmt.Sprint(v)
	}
	if v, ok := data["name"]; ok {
		recordAttributes.DataName = v.(string)
	}
	if v, ok := data["order"]; ok {
		recordAttributes.DataOrder = fmt.Sprint(v)
	}
	if v, ok := data["port"]; ok {
		recordAttributes.DataPort = fmt.Sprint(v)
	}
	if v, ok := data["preference"]; ok {
		recordAttributes.DataPreference = fmt.Sprint(v)
	}
	if v, ok := data["priority"]; ok {
		recordAttributes.DataPriority = fmt.Sprint(v)
	}
	if v, ok := data["proto"]; ok {
		recordAttributes.DataProto = v.(string)
	}
	if v, ok := data["protocol"]; ok {
		recordAttributes.DataProtocol = fmt.Sprint(v)
	}
	if v, ok := data["public_key"]; ok {
		recordAttributes.DataPublicKey = fmt.Sprint(v)
	}
	if v, ok := data["regex"]; ok {
		recordAttributes.DataRegex = v.(string)
	}
	if v, ok := data["replacement"]; ok {
		recordAttributes.DataReplacement = v.(string)
	}
	if v, ok := data["selector"]; ok {
		recordAttributes.DataSelector = fmt.Sprint(v)
	}
	if v, ok := data["service"]; ok {
		recordAttributes.DataService = v.(string)
	}
	if v, ok := data["target"]; ok {
		recordAttributes.DataTarget = v.(string)
	}
	if v, ok := data["tag"]; ok {
		recordAttributes.DataTag = v.(string)
	}
	if v, ok := data["type"]; ok {
		recordAttributes.DataType = fmt.Sprint(v)
	}
	if v, ok := data["usage"]; ok {
		recordAttributes.DataUsage = fmt.Sprint(v)
	}
	if v, ok := data["value"]; ok {
		recordAttributes.DataValue = v.(string)
	}
	if v, ok := data["weight"]; ok {
		recordAttributes.DataWeight = fmt.Sprint(v)
	}
	if v, ok := data["altitude"]; ok {
		recordAttributes.DataAltitude = fmt.Sprint(v)
	}
	if v, ok := data["lat_degrees"]; ok {
		recordAttributes.DataLatDegrees = fmt.Sprint(v)
	}
	if v, ok := data["lat_direction"]; ok {
		recordAttributes.DataLatDirection = fmt.Sprint(v)
	}
	if v, ok := data["lat_minutes"]; ok {
		recordAttributes.DataLatMinutes = fmt.Sprint(v)
	}
	if v, ok := data["lat_seconds"]; ok {
		recordAttributes.DataLatSeconds = fmt.Sprint(v)
	}
	if v, ok := data["long_degrees"]; ok {
		recordAttributes.DataLongDegrees = fmt.Sprint(v)
	}
	if v, ok := data["long_direction"]; ok {
		recordAttributes.DataLongDirection = fmt.Sprint(v)
	}
	if v, ok := data["long_minutes"]; ok {
		recordAttributes.DataLongMinutes = fmt.Sprint(v)
	}
	if v, ok := data["long_seconds"]; ok {
		recordAttributes.DataLongSeconds = fmt.Sprint(v)
	}
	if v, ok := data["precision_horz"]; ok {
		recordAttributes.DataPrecisionHorz = fmt.Sprint(v)
	}
	if v, ok := data["precision_vert"]; ok {
		recordAttributes.DataPrecisionVert = fmt.Sprint(v)
	}
	if v, ok := data["size"]; ok {
		recordAttributes.DataSize = fmt.Sprint(v)
	}

	r.Primary.Attributes = recordAttributes

	return r
}
