package cmd

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
)

var (
	// listOfString is an example representation of a key where the value is a
	// list of string values.
	//
	//   resource "example" "example" {
	//     attr = [ "b", "c", "d"]
	//   }
	listOfString = []interface{}{"b", "c", "d"}

	// configBlockOfStrings is an example of where a key is a "block" assignment
	// in HCL.
	//
	//   resource "example" "example" {
	//     attr = {
	//       c = "d"
	//       e = "f"
	//     }
	//   }
	configBlockOfStrings = map[string]interface{}{
		"c": "d",
		"e": "f",
	}

	cloudflareTestZoneID    = "0da42c8d2132a9ddaf714f9e7c920711"
	cloudflareTestAccountID = "f037e56e89293a057740de681ac9abbe"
)

func TestGenerate_writeAttrLine(t *testing.T) {
	tests := map[string]struct {
		key   string
		value interface{}
		want  string
	}{
		"value is string":           {key: "a", value: "b", want: fmt.Sprintf("a = %q\n", "b")},
		"value is int":              {key: "a", value: 1, want: "a = 1\n"},
		"value is float":            {key: "a", value: 1.0, want: "a = 1\n"},
		"value is bool":             {key: "a", value: true, want: "a = true\n"},
		"value is list of strings":  {key: "a", value: listOfString, want: "a = [ \"b\", \"c\", \"d\" ]\n"},
		"value is block of strings": {key: "a", value: configBlockOfStrings, want: "a = {\nc = \"d\"\ne = \"f\"\n}\n"},
		"value is nil":              {key: "a", value: nil, want: ""},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := writeAttrLine(tc.key, tc.value, false)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGenerate_ResourceNotSupported(t *testing.T) {
	_, output, err := executeCommandC(rootCmd, "generate", "--resource-type", "notreal")

	if assert.Nil(t, err) {
		assert.Contains(t, output, "\"notreal\" is not yet supported for automatic generation")
	}
}

func TestResourceGeneration(t *testing.T) {
	tests := map[string]struct {
		identiferType    string
		resourceType     string
		testdataFilename string
	}{
		"cloudflare access application simple (account)":    {identiferType: "account", resourceType: "cloudflare_access_application", testdataFilename: "cloudflare_access_application_simple_account"},
		"cloudflare access application with CORS (account)": {identiferType: "account", resourceType: "cloudflare_access_application", testdataFilename: "cloudflare_access_application_with_cors_account"},
		"cloudflare access IdP OAuth (account)":             {identiferType: "account", resourceType: "cloudflare_access_identity_provider", testdataFilename: "cloudflare_access_identity_provider_oauth_account"},
		"cloudflare access IdP OAuth (zone)":                {identiferType: "zone", resourceType: "cloudflare_access_identity_provider", testdataFilename: "cloudflare_access_identity_provider_oauth_zone"},
		"cloudflare access IdP OTP (account)":               {identiferType: "account", resourceType: "cloudflare_access_identity_provider", testdataFilename: "cloudflare_access_identity_provider_otp_account"},
		"cloudflare access IdP OTP (zone)":                  {identiferType: "zone", resourceType: "cloudflare_access_identity_provider", testdataFilename: "cloudflare_access_identity_provider_otp_zone"},
		"cloudflare access rule (account)":                  {identiferType: "account", resourceType: "cloudflare_access_rule", testdataFilename: "cloudflare_access_rule_account"},
		"cloudflare access rule (zone)":                     {identiferType: "zone", resourceType: "cloudflare_access_rule", testdataFilename: "cloudflare_access_rule_zone"},
		"cloudflare account member":                         {identiferType: "account", resourceType: "cloudflare_account_member", testdataFilename: "cloudflare_account_member"},
		"cloudflare argo tunnel":                            {identiferType: "account", resourceType: "cloudflare_argo_tunnel", testdataFilename: "cloudflare_argo_tunnel"},
		"cloudflare argo":                                   {identiferType: "zone", resourceType: "cloudflare_argo", testdataFilename: "cloudflare_argo"},
		"cloudflare BYO IP prefix":                          {identiferType: "account", resourceType: "cloudflare_byo_ip_prefix", testdataFilename: "cloudflare_byo_ip_prefix"},
		"cloudflare certificate pack":                       {identiferType: "zone", resourceType: "cloudflare_certificate_pack", testdataFilename: "cloudflare_certificate_pack"},
		"cloudflare custom hostname fallback origin":        {identiferType: "zone", resourceType: "cloudflare_custom_hostname_fallback_origin", testdataFilename: "cloudflare_custom_hostname_fallback_origin"},
		"cloudflare custom hostname":                        {identiferType: "zone", resourceType: "cloudflare_custom_hostname", testdataFilename: "cloudflare_custom_hostname"},
		"cloudflare custom pages (account)":                 {identiferType: "account", resourceType: "cloudflare_custom_pages", testdataFilename: "cloudflare_custom_pages_account"},
		"cloudflare custom pages (zone)":                    {identiferType: "zone", resourceType: "cloudflare_custom_pages", testdataFilename: "cloudflare_custom_pages_zone"},
		"cloudflare filter":                                 {identiferType: "zone", resourceType: "cloudflare_filter", testdataFilename: "cloudflare_filter"},
		"cloudflare firewall rule":                          {identiferType: "zone", resourceType: "cloudflare_firewall_rule", testdataFilename: "cloudflare_firewall_rule"},
		"cloudflare health check":                           {identiferType: "zone", resourceType: "cloudflare_healthcheck", testdataFilename: "cloudflare_healthcheck"},
		"cloudflare load balancer monitor":                  {identiferType: "account", resourceType: "cloudflare_load_balancer_monitor", testdataFilename: "cloudflare_load_balancer_monitor"},
		"cloudflare load balancer":                          {identiferType: "zone", resourceType: "cloudflare_load_balancer", testdataFilename: "cloudflare_load_balancer"},
		"cloudflare logpush jobs":                           {identiferType: "zone", resourceType: "cloudflare_logpush_job", testdataFilename: "cloudflare_logpush_job"},
		"cloudflare origin CA certificate":                  {identiferType: "zone", resourceType: "cloudflare_origin_ca_certificate", testdataFilename: "cloudflare_origin_ca_certificate"},
		"cloudflare page rule":                              {identiferType: "zone", resourceType: "cloudflare_page_rule", testdataFilename: "cloudflare_page_rule"},
		"cloudflare rate limit":                             {identiferType: "zone", resourceType: "cloudflare_rate_limit", testdataFilename: "cloudflare_rate_limit"},
		"cloudflare record CAA":                             {identiferType: "zone", resourceType: "cloudflare_record", testdataFilename: "cloudflare_record_caa"},
		"cloudflare record PTR":                             {identiferType: "zone", resourceType: "cloudflare_record", testdataFilename: "cloudflare_record_ptr"},
		"cloudflare record TXT SPF":                         {identiferType: "zone", resourceType: "cloudflare_record", testdataFilename: "cloudflare_record_txt_spf"},
		"cloudflare record simple":                          {identiferType: "zone", resourceType: "cloudflare_record", testdataFilename: "cloudflare_record"},
		"cloudflare record subdomain":                       {identiferType: "zone", resourceType: "cloudflare_record", testdataFilename: "cloudflare_record_subdomain"},
		"cloudflare ruleset":                                {identiferType: "zone", resourceType: "cloudflare_ruleset", testdataFilename: "cloudflare_ruleset_zone"},
		"cloudflare spectrum application":                   {identiferType: "zone", resourceType: "cloudflare_spectrum_application", testdataFilename: "cloudflare_spectrum_application"},
		"cloudflare WAF override":                           {identiferType: "zone", resourceType: "cloudflare_waf_override", testdataFilename: "cloudflare_waf_override"},
		"cloudflare waiting room":                           {identiferType: "zone", resourceType: "cloudflare_waiting_room", testdataFilename: "cloudflare_waiting_room"},
		"cloudflare worker route":                           {identiferType: "zone", resourceType: "cloudflare_worker_route", testdataFilename: "cloudflare_worker_route"},
		"cloudflare workers kv namespace":                   {identiferType: "account", resourceType: "cloudflare_workers_kv_namespace", testdataFilename: "cloudflare_workers_kv_namespace"},
		"cloudflare zone lockdown":                          {identiferType: "zone", resourceType: "cloudflare_zone_lockdown", testdataFilename: "cloudflare_zone_lockdown"},
		"cloudflare zone settings override":                 {identiferType: "zone", resourceType: "cloudflare_zone_settings_override", testdataFilename: "cloudflare_zone_settings_override"},

		// "cloudflare access group (account)": {identiferType: "account", resourceType: "cloudflare_access_group", testdataFilename: "cloudflare_access_group_account"},
		// "cloudflare access group (zone)":    {identiferType: "zone", resourceType: "cloudflare_access_group", testdataFilename: "cloudflare_access_group_zone"},
		// "cloudflare custom certificates":    {identiferType: "zone", resourceType: "cloudflare_custom_certificates", testdataFilename: "cloudflare_custom_certificates"},
		// "cloudflare custom SSL":             {identiferType: "zone", resourceType: "cloudflare_custom_ssl", testdataFilename: "cloudflare_custom_ssl"},
		// "cloudflare load balancer pool":     {identiferType: "account", resourceType: "cloudflare_load_balancer_pool", testdataFilename: "cloudflare_load_balancer_pool"},
		// "cloudflare worker cron trigger":    {identiferType: "zone", resourceType: "cloudflare_worker_cron_trigger", testdataFilename: "cloudflare_worker_cron_trigger"},
		// "cloudflare zone":                   {identiferType: "zone", resourceType: "cloudflare_zone", testdataFilename: "cloudflare_zone"},
	}

	for name, tc := range tests {

		t.Run(name, func(t *testing.T) {
			// Reset the environment variables used in test to ensure we don't
			// have both present at once.
			viper.Set("zone", "")
			viper.Set("account", "")

			r, err := recorder.New("../../../../testdata/cloudflare/" + tc.testdataFilename)
			if err != nil {
				log.Fatal(err)
			}
			defer r.Stop()

			r.AddFilter(func(i *cassette.Interaction) error {
				delete(i.Request.Headers, "X-Auth-Email")
				delete(i.Request.Headers, "X-Auth-Key")
				delete(i.Request.Headers, "Authorization")
				return nil
			})

			output := ""

			if tc.identiferType == "account" {
				viper.Set("account", cloudflareTestAccountID)
				api, _ = cloudflare.New(viper.GetString("key"), viper.GetString("email"), cloudflare.HTTPClient(
					&http.Client{
						Transport: r,
					},
				), cloudflare.UsingAccount(cloudflareTestAccountID))

				_, output, _ = executeCommandC(rootCmd, "generate", "--resource-type", tc.resourceType, "--account", cloudflareTestAccountID)

			} else {
				viper.Set("zone", cloudflareTestZoneID)
				api, _ = cloudflare.New(viper.GetString("key"), viper.GetString("email"), cloudflare.HTTPClient(
					&http.Client{
						Transport: r,
					},
				))

				_, output, _ = executeCommandC(rootCmd, "generate", "--resource-type", tc.resourceType, "--zone", cloudflareTestZoneID)

			}

			expected := testDataFile(tc.testdataFilename + ".tf")
			assert.Equal(t, strings.TrimRight(expected, "\n"), strings.TrimRight(output, "\n"))
		})
	}
}
