package cmd

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"strings"
	"time"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/spf13/cobra"

	"fmt"
)

// resourceImportStringFormats contains a mapping of the resource type to the
// composite ID that is compatible with performing an import.
var resourceImportStringFormats = map[string]string{
	"cloudflare_access_rule":           ":identifer_type/:identifer_value/:id",
	"cloudflare_account_member":        ":account_id/:id",
	"cloudflare_argo":                  ":zone_id/argo",
	"cloudflare_argo_tunnel":           ":account_id/:id",
	"cloudflare_byo_ip_prefix":         ":id",
	"cloudflare_certificate_pack":      ":zone_id/:id",
	"cloudflare_custom_pages":          ":identifer_type/:identifer_value/:id",
	"cloudflare_filter":                ":zone_id/:id",
	"cloudflare_firewall_rule":         ":zone_id/:id",
	"cloudflare_custom_hostname":       ":zone_id/:id",
	"cloudflare_custom_ssl":            ":zone_id/:id",
	"cloudflare_ip_list":               ":account_id/:id",
	"cloudflare_origin_ca_certificate": ":id",
	"cloudflare_page_rule":             ":zone_id/:id",
	"cloudflare_rate_limit":            ":zone_id/:id",
	"cloudflare_record":                ":zone_id/:id",
	"cloudflare_spectrum_application":  ":zone_id/:id",
	"cloudflare_waf_override":          ":zone_id/:id",
	"cloudflare_waf_package":           ":zone_id/:id",
	"cloudflare_waf_group":             ":zone_id/:id",
	"cloudflare_worker_route":          ":zone_id/:id",
	"cloudflare_workers_kv_namespace":  ":id",
	"cloudflare_zone":                  ":id",
	"cloudflare_zone_lockdown":         ":zone_id/:id",
}

func init() {
	rootCmd.AddCommand(importCommand)
}

var importCommand = &cobra.Command{
	Use:    "import",
	Short:  "Output `terraform import` compatible commands in order to import resources into state",
	Run:    runImport(),
	PreRun: sharedPreRun,
}

func runImport() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		var jsonStructData []interface{}
		switch resourceType {
		case "cloudflare_access_rule":
			if accountID != "" {
				jsonPayload, err := api.ListAccountAccessRules(context.Background(), accountID, cloudflare.AccessRule{}, 1)
				if err != nil {
					log.Fatal(err)
				}

				m, _ := json.Marshal(jsonPayload.Result)
				json.Unmarshal(m, &jsonStructData)
			} else {
				jsonPayload, err := api.ListZoneAccessRules(context.Background(), zoneID, cloudflare.AccessRule{}, 1)
				if err != nil {
					log.Fatal(err)
				}

				m, _ := json.Marshal(jsonPayload.Result)
				json.Unmarshal(m, &jsonStructData)
			}
		case "cloudflare_account_member":
			jsonPayload, _, err := api.AccountMembers(context.Background(), accountID, cloudflare.PaginationOptions{})
			if err != nil {
				log.Fatal(err)
			}
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_argo":
			jsonPayload := []cloudflare.ArgoFeatureSetting{{
				ID: fmt.Sprintf("%x", md5.Sum([]byte(time.Now().String()))),
			}}

			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_argo_tunnel":
			jsonPayload, err := api.ArgoTunnels(context.Background(), accountID)
			if err != nil {
				log.Fatal(err)
			}
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_byo_ip_prefix":
			jsonPayload, err := api.ListPrefixes(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_certificate_pack":
			jsonPayload, err := api.ListCertificatePacks(context.Background(), zoneID)
			if err != nil {
				log.Fatal(err)
			}
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_custom_pages":
			if accountID != "" {
				jsonPayload, err := api.CustomPages(context.Background(), &cloudflare.CustomPageOptions{AccountID: accountID})
				if err != nil {
					log.Fatal(err)
				}

				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)

			} else {
				jsonPayload, err := api.CustomPages(context.Background(), &cloudflare.CustomPageOptions{ZoneID: zoneID})
				if err != nil {
					log.Fatal(err)
				}

				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			}
		case "cloudflare_filter":
			jsonPayload, err := api.Filters(context.Background(), zoneID, cloudflare.PaginationOptions{})
			if err != nil {
				log.Fatal(err)
			}
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_firewall_rule":
			jsonPayload, err := api.FirewallRules(context.Background(), zoneID, cloudflare.PaginationOptions{})
			if err != nil {
				log.Fatal(err)
			}
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_custom_hostname":
			jsonPayload, _, err := api.CustomHostnames(context.Background(), zoneID, 1, cloudflare.CustomHostname{})
			if err != nil {
				log.Fatal(err)
			}
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_custom_ssl":
			jsonPayload, err := api.ListSSL(context.Background(), zoneID)
			if err != nil {
				log.Fatal(err)
			}

			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_ip_list":
			jsonPayload, err := api.ListIPLists(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_logpush_job":
			jsonPayload, err := api.LogpushJobs(context.Background(), zoneID)
			if err != nil {
				log.Fatal(err)
			}
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_origin_ca_certificate":
			jsonPayload, err := api.OriginCertificates(context.Background(), cloudflare.OriginCACertificateListOptions{ZoneID: zoneID})
			if err != nil {
				log.Fatal(err)
			}

			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_page_rule":
			jsonPayload, err := api.ListPageRules(context.Background(), zoneID)
			if err != nil {
				log.Fatal(err)
			}

			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_rate_limit":
			jsonPayload, _, err := api.ListRateLimits(context.Background(), zoneID, cloudflare.PaginationOptions{})
			if err != nil {
				log.Fatal(err)
			}

			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_record":
			jsonPayload, err := api.DNSRecords(context.Background(), zoneID, cloudflare.DNSRecord{})
			if err != nil {
				log.Fatal(err)
			}
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_spectrum_application":
			jsonPayload, err := api.SpectrumApplications(context.Background(), zoneID)
			if err != nil {
				log.Fatal(err)
			}

			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_waf_override":
			jsonPayload, err := api.ListWAFOverrides(context.Background(), zoneID)
			if err != nil {
				log.Fatal(err)
			}

			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_waf_package":
			jsonPayload, err := api.ListWAFPackages(context.Background(), zoneID)
			if err != nil {
				log.Fatal(err)
			}
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_waf_group":
			jsonPayload, err := api.ListWAFGroups(context.Background(), zoneID, packageID)
			if err != nil {
				log.Fatal(err)
			}
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_workers_kv_namespace":
			jsonPayload, err := api.ListWorkersKVNamespaces(context.Background())
			if err != nil {
				log.Fatal(err)
			}

			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_worker_route":
			jsonPayload, err := api.ListWorkerRoutes(context.Background(), zoneID)
			if err != nil {
				log.Fatal(err)
			}

			m, _ := json.Marshal(jsonPayload.Routes)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_zone":
			jsonPayload, err := api.ListZones(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
		case "cloudflare_zone_lockdown":
			jsonPayload, err := api.ListZoneLockdowns(context.Background(), zoneID, 1)
			if err != nil {
				log.Fatal(err)
			}

			m, _ := json.Marshal(jsonPayload.Result)
			json.Unmarshal(m, &jsonStructData)
		default:
			fmt.Fprintf(cmd.OutOrStdout(), "%q is not yet supported for state import", resourceType)
			return
		}

		for _, data := range jsonStructData {
			fmt.Fprint(cmd.OutOrStdout(), buildCompositeID(resourceType, data.(map[string]interface{})["id"].(string)))
		}
	}
}

// buildCompositeID takes the resourceType and resourceID in order to lookup the
// resource type import string and then return a suitable composite value that
// is compatible with `terraform import`.
func buildCompositeID(resourceType, resourceID string) string {
	if _, ok := resourceImportStringFormats[resourceType]; !ok {
		log.Fatalf("%s does not have an import format defined", resourceType)
	}

	var identiferType string
	var identiferValue string

	if accountID != "" {
		identiferType = "account"
		identiferValue = accountID
	} else {
		identiferType = "zone"
		identiferValue = zoneID
	}

	s := ""
	s += fmt.Sprintf("%s %s.%s_%s %s", terraformImportCmdPrefix, resourceType, terraformResourceNamePrefix, resourceID, resourceImportStringFormats[resourceType])
	replacer := strings.NewReplacer(
		":identifer_type", identiferType,
		":identifer_value", identiferValue,
		":zone_id", zoneID,
		":account_id", accountID,
		":id", resourceID,
	)
	s += "\n"

	return replacer.Replace(s)
}
