package cmd

import (
	"context"
	"encoding/json"
	"strings"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/spf13/cobra"

	"fmt"
)

// resourceImportStringFormats contains a mapping of the resource type to the
// composite ID that is compatible with performing an import.
var resourceImportStringFormats = map[string]string{
	"cloudflare_account_member":   ":account_id/:id",
	"cloudflare_argo_tunnel":      ":account_id/:id",
	"cloudflare_byo_ip_prefix":    ":id",
	"cloudflare_certificate_pack": ":zone_id/:id",
	"cloudflare_filter":           ":zone_id/:id",
	"cloudflare_firewall_rule":    ":zone_id/:id",
	"cloudflare_custom_hostname":  ":zone_id/:id",
	"cloudflare_ip_list":          ":account_id/:id",
	"cloudflare_record":           ":zone_id/:id",
	"cloudflare_worker_route":     ":zone_id/:id",
	"cloudflare_zone":             ":zone_id",
}

func init() {
	rootCmd.AddCommand(importCommand)
}

var importCommand = &cobra.Command{
	Use:   "import",
	Short: "Output `terraform import` compatible commands in order to import resources into state",
	Run:   runImport(),
}

func runImport() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		var jsonStructData []interface{}
		switch resourceType {
		case "cloudflare_access_identity_provider":
			if accountID != "" {
				jsonPayload, err := api.AccessIdentityProviders(context.Background(), accountID)
				if err != nil {
					log.Fatal(err)
				}

				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			} else {
				jsonPayload, err := api.ZoneLevelAccessIdentityProviders(context.Background(), zoneID)
				if err != nil {
					log.Fatal(err)
				}

				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			}
		case "cloudflare_access_service_token":
			if accountID != "" {
				jsonPayload, _, err := api.AccessServiceTokens(context.Background(), accountID)
				if err != nil {
					log.Fatal(err)
				}

				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			} else {
				jsonPayload, _, err := api.ZoneLevelAccessServiceTokens(context.Background(), zoneID)
				if err != nil {
					log.Fatal(err)
				}

				m, _ := json.Marshal(jsonPayload)
				json.Unmarshal(m, &jsonStructData)
			}
		case "cloudflare_access_mutual_tls_certificate":
			jsonPayload, err := api.AccessMutualTLSCertificates(context.Background(), accountID)
			if err != nil {
				log.Fatal(err)
			}
			m, _ := json.Marshal(jsonPayload)
			json.Unmarshal(m, &jsonStructData)
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
		case "cloudflare_record":
			jsonPayload, err := api.DNSRecords(context.Background(), zoneID, cloudflare.DNSRecord{})
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

	s := ""
	s += fmt.Sprintf("%s %s.%s_%s %s", terraformImportCmdPrefix, resourceType, terraformResourceNamePrefix, resourceID, resourceImportStringFormats[resourceType])
	replacer := strings.NewReplacer(
		":zone_id", zoneID,
		":account_id", accountID,
		":id", resourceID,
	)
	s += "\n"

	return replacer.Replace(s)
}
