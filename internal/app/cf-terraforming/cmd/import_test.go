package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/option"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
)

func TestResourceImportV5(t *testing.T) {
	tests := map[string]struct {
		identiferType    string
		resourceType     string
		testdataFilename string
		cliFlags         string
	}{
		"cloudflare account":                                       {identiferType: "account", resourceType: "cloudflare_account", testdataFilename: "cloudflare_account"},
		"cloudflare address map":                                   {identiferType: "account", resourceType: "cloudflare_address_map", testdataFilename: "cloudflare_address_map"},
		"cloudflare account member":                                {identiferType: "account", resourceType: "cloudflare_account_member", testdataFilename: "cloudflare_account_member"},
		"cloudflare api shield operation":                          {identiferType: "zone", resourceType: "cloudflare_api_shield_operation", testdataFilename: "cloudflare_api_shield_operation"},
		"cloudflare bot management":                                {identiferType: "zone", resourceType: "cloudflare_bot_management", testdataFilename: "cloudflare_bot_management"},
		"cloudflare certificate pack":                              {identiferType: "zone", resourceType: "cloudflare_certificate_pack", testdataFilename: "cloudflare_certificate_pack"},
		"cloudflare custom hostname fallback origin":               {identiferType: "zone", resourceType: "cloudflare_custom_hostname_fallback_origin", testdataFilename: "cloudflare_custom_hostname_fallback_origin"},
		"cloudflare custom hostname":                               {identiferType: "zone", resourceType: "cloudflare_custom_hostname", testdataFilename: "cloudflare_custom_hostname"},
		"cloudflare email routing address":                         {identiferType: "account", resourceType: "cloudflare_email_routing_address", testdataFilename: "cloudflare_email_routing_address"},
		"cloudflare email routing catch all":                       {identiferType: "zone", resourceType: "cloudflare_email_routing_catch_all", testdataFilename: "cloudflare_email_routing_catch_all"},
		"cloudflare email routing dns":                             {identiferType: "zone", resourceType: "cloudflare_email_routing_dns", testdataFilename: "cloudflare_email_routing_dns"},
		"cloudflare email routing rule":                            {identiferType: "zone", resourceType: "cloudflare_email_routing_rule", testdataFilename: "cloudflare_email_routing_rule"},
		"cloudflare email routing settings":                        {identiferType: "zone", resourceType: "cloudflare_email_routing_settings", testdataFilename: "cloudflare_email_routing_settings"},
		"cloudflare email security block sender":                   {identiferType: "account", resourceType: "cloudflare_email_security_block_sender", testdataFilename: "cloudflare_email_security_block_sender"},
		"cloudflare email security trusted domains":                {identiferType: "account", resourceType: "cloudflare_email_security_trusted_domains", testdataFilename: "cloudflare_email_security_trusted_domains"},
		"cloudflare email security impersonation registry":         {identiferType: "account", resourceType: "cloudflare_email_security_impersonation_registry", testdataFilename: "cloudflare_email_security_impersonation_registry"},
		"cloudflare filter":                                        {identiferType: "zone", resourceType: "cloudflare_filter", testdataFilename: "cloudflare_filter"},
		"cloudflare health check":                                  {identiferType: "zone", resourceType: "cloudflare_healthcheck", testdataFilename: "cloudflare_healthcheck"},
		"cloudflare hostname tls setting":                          {identiferType: "zone", resourceType: "cloudflare_hostname_tls_setting", testdataFilename: "cloudflare_hostname_tls_setting", cliFlags: "cloudflare_hostname_tls_setting=ciphers,min_tls_version"},
		"cloudflare keyless certificate":                           {identiferType: "zone", resourceType: "cloudflare_keyless_certificate", testdataFilename: "cloudflare_keyless_certificate"},
		"cloudflare mtls certificate":                              {identiferType: "account", resourceType: "cloudflare_mtls_certificate", testdataFilename: "cloudflare_mtls_certificate"},
		"cloudflare load balancer":                                 {identiferType: "zone", resourceType: "cloudflare_load_balancer", testdataFilename: "cloudflare_load_balancer"},
		"cloudflare load balancer monitor":                         {identiferType: "account", resourceType: "cloudflare_load_balancer_monitor", testdataFilename: "cloudflare_load_balancer_monitor"},
		"cloudflare load balancer pool":                            {identiferType: "account", resourceType: "cloudflare_load_balancer_pool", testdataFilename: "cloudflare_load_balancer_pool"},
		"cloudflare managed transforms":                            {identiferType: "zone", resourceType: "cloudflare_managed_transforms", testdataFilename: "cloudflare_managed_transforms"},
		"cloudflare origin ca certificate":                         {identiferType: "zone", resourceType: "cloudflare_origin_ca_certificate", testdataFilename: "cloudflare_origin_ca_certificate"},
		"cloudflare d1 database":                                   {identiferType: "account", resourceType: "cloudflare_d1_database", testdataFilename: "cloudflare_d1_database"},
		"cloudflare dns firewall":                                  {identiferType: "account", resourceType: "cloudflare_dns_firewall", testdataFilename: "cloudflare_dns_firewall"},
		"cloudflare dns record simple":                             {identiferType: "zone", resourceType: "cloudflare_dns_record", testdataFilename: "cloudflare_dns_record"},
		"cloudflare dns zone transfers acl":                        {identiferType: "account", resourceType: "cloudflare_dns_zone_transfers_acl", testdataFilename: "cloudflare_dns_zone_transfers_acl"},
		"cloudflare dns zone transfers incoming":                   {identiferType: "zone", resourceType: "cloudflare_dns_zone_transfers_incoming", testdataFilename: "cloudflare_dns_zone_transfers_incoming"},
		"cloudflare dns zone transfers outgoing":                   {identiferType: "zone", resourceType: "cloudflare_dns_zone_transfers_outgoing", testdataFilename: "cloudflare_dns_zone_transfers_outgoing"},
		"cloudflare dns zone transfers peer":                       {identiferType: "account", resourceType: "cloudflare_dns_zone_transfers_peer", testdataFilename: "cloudflare_dns_zone_transfers_peer"},
		"cloudflare dns zone transfers tsig":                       {identiferType: "account", resourceType: "cloudflare_dns_zone_transfers_tsig", testdataFilename: "cloudflare_dns_zone_transfers_tsig"},
		"cloudflare list":                                          {identiferType: "account", resourceType: "cloudflare_list", testdataFilename: "cloudflare_list"},
		"cloudflare list item":                                     {identiferType: "account", resourceType: "cloudflare_list_item", testdataFilename: "cloudflare_list_item", cliFlags: "cloudflare_list_item=2a4b8b2017aa4b3cb9e1151b52c81d22"},
		"cloudflare logpush job":                                   {identiferType: "account", resourceType: "cloudflare_logpush_job", testdataFilename: "cloudflare_logpush_job"},
		"cloudflare notification policy":                           {identiferType: "account", resourceType: "cloudflare_notification_policy", testdataFilename: "cloudflare_notification_policy"},
		"cloudflare notification policy webhooks":                  {identiferType: "account", resourceType: "cloudflare_notification_policy_webhooks", testdataFilename: "cloudflare_notification_policy_webhooks"},
		"cloudflare pages domain":                                  {identiferType: "account", resourceType: "cloudflare_pages_domain", testdataFilename: "cloudflare_pages_domain", cliFlags: "cloudflare_pages_domain=ykfjmcgpfs"},
		"cloudflare pages project":                                 {identiferType: "account", resourceType: "cloudflare_pages_project", testdataFilename: "cloudflare_pages_project"},
		"cloudflare page shield policy":                            {identiferType: "zone", resourceType: "cloudflare_page_shield_policy", testdataFilename: "cloudflare_page_shield_policy"},
		"cloudflare r2 bucket":                                     {identiferType: "account", resourceType: "cloudflare_r2_bucket", testdataFilename: "cloudflare_r2_bucket"},
		"cloudflare r2 managed domain":                             {identiferType: "account", resourceType: "cloudflare_r2_managed_domain", testdataFilename: "cloudflare_r2_managed_domain", cliFlags: "cloudflare_r2_managed_domain=jb-test-bucket,bnfywlzwpt"},
		"cloudflare r2 custom domain":                              {identiferType: "account", resourceType: "cloudflare_r2_custom_domain", testdataFilename: "cloudflare_r2_custom_domain", cliFlags: "cloudflare_r2_custom_domain=jb-test-bucket,bnfywlzwpt"},
		"cloudflare page rule":                                     {identiferType: "zone", resourceType: "cloudflare_page_rule", testdataFilename: "cloudflare_page_rule"},
		"cloudflare rate limit":                                    {identiferType: "zone", resourceType: "cloudflare_rate_limit", testdataFilename: "cloudflare_rate_limit"},
		"cloudflare ruleset (ddos_l7)":                             {identiferType: "zone", resourceType: "cloudflare_ruleset", testdataFilename: "cloudflare_ruleset_zone_ddos_l7"},
		"cloudflare ruleset (http_log_custom_fields)":              {identiferType: "zone", resourceType: "cloudflare_ruleset", testdataFilename: "cloudflare_ruleset_zone_http_log_custom_fields"},
		"cloudflare ruleset (http_ratelimit)":                      {identiferType: "zone", resourceType: "cloudflare_ruleset", testdataFilename: "cloudflare_ruleset_zone_http_ratelimit"},
		"cloudflare ruleset (http_request_cache_settings)":         {identiferType: "zone", resourceType: "cloudflare_ruleset", testdataFilename: "cloudflare_ruleset_http_request_cache_settings"},
		"cloudflare ruleset (http_request_firewall_custom)":        {identiferType: "zone", resourceType: "cloudflare_ruleset", testdataFilename: "cloudflare_ruleset_zone_http_request_firewall_custom"},
		"cloudflare ruleset (http_request_firewall_managed)":       {identiferType: "zone", resourceType: "cloudflare_ruleset", testdataFilename: "cloudflare_ruleset_zone_http_request_firewall_managed"},
		"cloudflare ruleset (http_request_late_transform)":         {identiferType: "zone", resourceType: "cloudflare_ruleset", testdataFilename: "cloudflare_ruleset_zone_http_request_late_transform"},
		"cloudflare ruleset (http_request_sanitize)":               {identiferType: "zone", resourceType: "cloudflare_ruleset", testdataFilename: "cloudflare_ruleset_zone_http_request_sanitize"},
		"cloudflare ruleset (no configuration)":                    {identiferType: "zone", resourceType: "cloudflare_ruleset", testdataFilename: "cloudflare_ruleset_zone_no_configuration"},
		"cloudflare ruleset (override remapping = disabled)":       {identiferType: "zone", resourceType: "cloudflare_ruleset", testdataFilename: "cloudflare_ruleset_override_remapping_disabled"},
		"cloudflare ruleset (override remapping = enabled)":        {identiferType: "zone", resourceType: "cloudflare_ruleset", testdataFilename: "cloudflare_ruleset_override_remapping_enabled"},
		"cloudflare ruleset (rewrite to empty query string)":       {identiferType: "zone", resourceType: "cloudflare_ruleset", testdataFilename: "cloudflare_ruleset_zone_rewrite_to_empty_query_parameter"},
		"cloudflare ruleset":                                       {identiferType: "zone", resourceType: "cloudflare_ruleset", testdataFilename: "cloudflare_ruleset"},
		"cloudflare spectrum application":                          {identiferType: "zone", resourceType: "cloudflare_spectrum_application", testdataFilename: "cloudflare_spectrum_application"},
		"cloudflare tiered cache":                                  {identiferType: "zone", resourceType: "cloudflare_tiered_cache", testdataFilename: "cloudflare_tiered_cache"},
		"cloudflare regional hostnames":                            {identiferType: "zone", resourceType: "cloudflare_regional_hostname", testdataFilename: "cloudflare_regional_hostname"},
		"cloudflare regional tiered cache":                         {identiferType: "zone", resourceType: "cloudflare_regional_tiered_cache", testdataFilename: "cloudflare_regional_tiered_cache"},
		"cloudflare total tls":                                     {identiferType: "zone", resourceType: "cloudflare_total_tls", testdataFilename: "cloudflare_total_tls"},
		"cloudflare turnstile widget no domains":                   {identiferType: "account", resourceType: "cloudflare_turnstile_widget", testdataFilename: "cloudflare_turnstile_widget_no_domains"},
		"cloudflare url normalization settings":                    {identiferType: "zone", resourceType: "cloudflare_url_normalization_settings", testdataFilename: "cloudflare_url_normalization_settings"},
		"cloudflare waiting room event":                            {identiferType: "zone", resourceType: "cloudflare_waiting_room_event", testdataFilename: "cloudflare_waiting_room_event", cliFlags: "cloudflare_waiting_room_event=e7f9e4c190ea8d6c66cab32ac110f39a"},
		"cloudflare waiting room rules":                            {identiferType: "zone", resourceType: "cloudflare_waiting_room_rules", testdataFilename: "cloudflare_waiting_room_rules", cliFlags: "cloudflare_waiting_room_rules=8bbd1b13450f6c63ab6ab4e08a63762d"},
		"cloudflare web3 hostname":                                 {identiferType: "zone", resourceType: "cloudflare_web3_hostname", testdataFilename: "cloudflare_web3_hostname"},
		"cloudflare zone lockdown":                                 {identiferType: "zone", resourceType: "cloudflare_zone_lockdown", testdataFilename: "cloudflare_zone_lockdown"},
		"cloudflare queue":                                         {identiferType: "account", resourceType: "cloudflare_queue", testdataFilename: "cloudflare_queue"},
		"cloudflare web analytics site":                            {identiferType: "account", resourceType: "cloudflare_web_analytics_site", testdataFilename: "cloudflare_web_analytics_site"},
		"cloudflare web analytics rule":                            {identiferType: "account", resourceType: "cloudflare_web_analytics_rule", testdataFilename: "cloudflare_web_analytics_rule", cliFlags: "cloudflare_web_analytics_rule=2fa89d8f-35f7-49ef-87d3-f24e866a5d5e"},
		"cloudflare waiting room":                                  {identiferType: "zone", resourceType: "cloudflare_waiting_room", testdataFilename: "cloudflare_waiting_room"},
		"cloudflare waiting room settings":                         {identiferType: "zone", resourceType: "cloudflare_waiting_room_settings", testdataFilename: "cloudflare_waiting_room_settings"},
		"cloudflare workers custom domain":                         {identiferType: "account", resourceType: "cloudflare_workers_custom_domain", testdataFilename: "cloudflare_workers_custom_domain"},
		"cloudflare workers kv namespace":                          {identiferType: "account", resourceType: "cloudflare_workers_kv_namespace", testdataFilename: "cloudflare_workers_kv_namespace"},
		"cloudflare workers for platforms dispatch namespace":      {identiferType: "account", resourceType: "cloudflare_workers_for_platforms_dispatch_namespace", testdataFilename: "cloudflare_workers_for_platforms_dispatch_namespace"},
		"cloudflare zero trust access application":                 {identiferType: "account", resourceType: "cloudflare_zero_trust_access_application", testdataFilename: "cloudflare_zero_trust_access_application"},
		"cloudflare zero trust access custom page":                 {identiferType: "account", resourceType: "cloudflare_zero_trust_access_custom_page", testdataFilename: "cloudflare_zero_trust_access_custom_page"},
		"cloudflare zero trust access group":                       {identiferType: "account", resourceType: "cloudflare_zero_trust_access_group", testdataFilename: "cloudflare_zero_trust_access_group"},
		"cloudflare zero trust access identity provider":           {identiferType: "zone", resourceType: "cloudflare_zero_trust_access_identity_provider", testdataFilename: "cloudflare_zero_trust_access_identity_provider"},
		"cloudflare zero trust access infrastructure target":       {identiferType: "account", resourceType: "cloudflare_zero_trust_access_infrastructure_target", testdataFilename: "cloudflare_zero_trust_access_infrastructure_target"},
		"cloudflare zero trust access key configuration":           {identiferType: "account", resourceType: "cloudflare_zero_trust_access_key_configuration", testdataFilename: "cloudflare_zero_trust_access_key_configuration"},
		"cloudflare zero trust access policy":                      {identiferType: "account", resourceType: "cloudflare_zero_trust_access_policy", testdataFilename: "cloudflare_zero_trust_access_policy"},
		"cloudflare zero trust access service token":               {identiferType: "account", resourceType: "cloudflare_zero_trust_access_service_token", testdataFilename: "cloudflare_zero_trust_access_service_token"},
		"cloudflare zero trust access tag":                         {identiferType: "account", resourceType: "cloudflare_zero_trust_access_tag", testdataFilename: "cloudflare_zero_trust_access_tag"},
		"cloudflare zero trust access short lived certificate":     {identiferType: "account", resourceType: "cloudflare_zero_trust_access_short_lived_certificate", testdataFilename: "cloudflare_zero_trust_access_short_lived_certificate"},
		"cloudflare zero trust risk scoring integration":           {identiferType: "account", resourceType: "cloudflare_zero_trust_risk_scoring_integration", testdataFilename: "cloudflare_zero_trust_risk_scoring_integration"},
		"cloudflare zero trust dex test":                           {identiferType: "account", resourceType: "cloudflare_zero_trust_dex_test", testdataFilename: "cloudflare_zero_trust_dex_test"},
		"cloudflare zero trust device custom profile":              {identiferType: "account", resourceType: "cloudflare_zero_trust_device_custom_profile", testdataFilename: "cloudflare_zero_trust_device_custom_profile"},
		"cloudflare zero trust device posture rule":                {identiferType: "account", resourceType: "cloudflare_zero_trust_device_posture_rule", testdataFilename: "cloudflare_zero_trust_device_posture_rule"},
		"cloudflare zero trust device posture integration":         {identiferType: "account", resourceType: "cloudflare_zero_trust_device_posture_integration", testdataFilename: "cloudflare_zero_trust_device_posture_integration"},
		"cloudflare zero trust device managed networks":            {identiferType: "account", resourceType: "cloudflare_zero_trust_device_managed_networks", testdataFilename: "cloudflare_zero_trust_device_managed_networks"},
		"cloudflare zero trust device default profile":             {identiferType: "account", resourceType: "cloudflare_zero_trust_device_default_profile", testdataFilename: "cloudflare_zero_trust_device_default_profile"},
		"cloudflare zero trust dlp predefined profile":             {identiferType: "account", resourceType: "cloudflare_zero_trust_dlp_predefined_profile", testdataFilename: "cloudflare_zero_trust_dlp_predefined_profile", cliFlags: "cloudflare_zero_trust_dlp_predefined_profile=c8932cc4-3312-4152-8041-f3f257122dc4,56a8c060-01bb-4f89-ba1e-3ad42770a342"},
		"cloudflare zero trust dns location":                       {identiferType: "account", resourceType: "cloudflare_zero_trust_dns_location", testdataFilename: "cloudflare_zero_trust_dns_location"},
		"cloudflare zero trust gateway certificate":                {identiferType: "account", resourceType: "cloudflare_zero_trust_gateway_certificate", testdataFilename: "cloudflare_zero_trust_gateway_certificate"},
		"cloudflare zero trust gateway policy":                     {identiferType: "account", resourceType: "cloudflare_zero_trust_gateway_policy", testdataFilename: "cloudflare_zero_trust_gateway_policy"},
		"cloudflare zero trust gateway proxy endpoint":             {identiferType: "account", resourceType: "cloudflare_zero_trust_gateway_proxy_endpoint", testdataFilename: "cloudflare_zero_trust_gateway_proxy_endpoint"},
		"cloudflare zero trust list":                               {identiferType: "account", resourceType: "cloudflare_zero_trust_list", testdataFilename: "cloudflare_zero_trust_list"},
		"cloudflare zero trust gateway settings":                   {identiferType: "account", resourceType: "cloudflare_zero_trust_gateway_settings", testdataFilename: "cloudflare_zero_trust_gateway_settings"},
		"cloudflare zero trust organization":                       {identiferType: "account", resourceType: "cloudflare_zero_trust_organization", testdataFilename: "cloudflare_zero_trust_organization"},
		"cloudflare zero trust tunnel cloudflared":                 {identiferType: "account", resourceType: "cloudflare_zero_trust_tunnel_cloudflared", testdataFilename: "cloudflare_zero_trust_tunnel_cloudflared"},
		"cloudflare zero trust tunnel cloudflared route":           {identiferType: "account", resourceType: "cloudflare_zero_trust_tunnel_cloudflared_route", testdataFilename: "cloudflare_zero_trust_tunnel_cloudflared_route"},
		"cloudflare zero trust tunnel cloudflared virtual network": {identiferType: "account", resourceType: "cloudflare_zero_trust_tunnel_cloudflared_virtual_network", testdataFilename: "cloudflare_zero_trust_tunnel_cloudflared_virtual_network"},
		"cloudflare zone":                                          {identiferType: "zone", resourceType: "cloudflare_zone", testdataFilename: "cloudflare_zone"},
		"cloudflare zone dnssec":                                   {identiferType: "zone", resourceType: "cloudflare_zone_dnssec", testdataFilename: "cloudflare_zone_dnssec"},
		"cloudflare zone setting":                                  {identiferType: "zone", resourceType: "cloudflare_zone_setting", testdataFilename: "cloudflare_zone_setting", cliFlags: "cloudflare_zone_setting=always_online,cache_level"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Reset the environment variables used in test to ensure we don't
			// have both present at once.
			viper.Set("zone", "")
			viper.Set("account", "")

			var r *recorder.Recorder
			var err error
			r, err = recorder.New("../../../../testdata/cloudflare/v5/" + tc.testdataFilename)

			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				err := r.Stop()
				if err != nil {
					log.Fatal(err)
				}
			}()

			r.AddFilter(func(i *cassette.Interaction) error {
				// Sensitive HTTP headers
				delete(i.Request.Headers, "X-Auth-Email")
				delete(i.Request.Headers, "X-Auth-Key")
				delete(i.Request.Headers, "Authorization")

				// HTTP request headers that we don't need to assert against
				delete(i.Request.Headers, "User-Agent")

				// HTTP response headers that we don't need to assert against
				delete(i.Response.Headers, "Cf-Cache-Status")
				delete(i.Response.Headers, "Cf-Ray")
				delete(i.Response.Headers, "Date")
				delete(i.Response.Headers, "Server")
				delete(i.Response.Headers, "Set-Cookie")
				delete(i.Response.Headers, "X-Envoy-Upstream-Service-Time")

				if os.Getenv("CLOUDFLARE_DOMAIN") != "" {
					i.Response.Body = strings.ReplaceAll(i.Response.Body, os.Getenv("CLOUDFLARE_DOMAIN"), "example.com")
				}

				return nil
			})

			output := ""
			api = cloudflare.NewClient(option.WithHTTPClient(
				&http.Client{
					Transport: r,
				},
			))
			if tc.identiferType == "account" {
				viper.Set("account", cloudflareTestAccountID)
				if tc.cliFlags != "" {
					output, err = executeCommandC(
						rootCmd,
						"import",
						"--resource-type",
						tc.resourceType,
						"--account",
						cloudflareTestAccountID,
						"--resource-id",
						tc.cliFlags,
					)
				} else {
					output, err = executeCommandC(rootCmd, "import", "--resource-type", tc.resourceType, "--account", cloudflareTestAccountID)
				}
			} else {
				viper.Set("zone", cloudflareTestZoneID)
				if tc.cliFlags != "" {
					output, err = executeCommandC(
						rootCmd,
						"import",
						"--resource-type",
						tc.resourceType,
						"--zone",
						cloudflareTestZoneID,
						"--resource-id",
						tc.cliFlags,
					)
				} else {
					output, err = executeCommandC(rootCmd, "import", "--resource-type", tc.resourceType, "--zone", cloudflareTestZoneID)
				}
			}
			assert.NotEmpty(t, output, fmt.Sprintf("should have output. But got %s", output))
			assert.NoError(t, err, fmt.Sprintf("should not have an error. But got %+v", err))
			assert.True(t, strings.Contains(output, "terraform import"), fmt.Sprintf("should match %s", output))
		})
	}
}
