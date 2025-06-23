package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"sort"
	"strings"

	"github.com/cloudflare/cloudflare-go/v4"

	"github.com/sirupsen/logrus"

	"github.com/tidwall/gjson"
)

func processCustomCasesV5(response *[]interface{}, resourceType string, pathParam string) {
	resourceCount := len(*response)
	switch resourceType {
	case "cloudflare_managed_transforms":
		// remap email and role_ids into the right structure and remove policies
		for i := 0; i < resourceCount; i++ {
			for j := range (*response)[i].(map[string]interface{})["managed_request_headers"].([]interface{}) {
				delete((*response)[i].(map[string]interface{})["managed_request_headers"].([]interface{})[j].(map[string]interface{}), "has_conflict")
			}
			for j := range (*response)[i].(map[string]interface{})["managed_response_headers"].([]interface{}) {
				delete((*response)[i].(map[string]interface{})["managed_response_headers"].([]interface{})[j].(map[string]interface{}), "has_conflict")
			}
		}
	case "cloudflare_r2_bucket":
		denestResponses(response, resourceCount, "buckets")
	case "cloudflare_account_member":
		// remap email and role_ids into the right structure and remove policies
		for i := 0; i < resourceCount; i++ {
			delete((*response)[i].(map[string]interface{}), "policies")
			(*response)[i].(map[string]interface{})["email"] = (*response)[i].(map[string]interface{})["user"].(map[string]interface{})["email"]
			roleIDs := []string{}
			for _, role := range (*response)[i].(map[string]interface{})["roles"].([]interface{}) {
				roleIDs = append(roleIDs, role.(map[string]interface{})["id"].(string))
			}
			(*response)[i].(map[string]interface{})["roles"] = roleIDs
		}
	case "cloudflare_content_scanning_expression":
		// wrap the response in 'body' for tf
		for i := 0; i < resourceCount; i++ {
			payload := (*response)[i].(map[string]interface{})["payload"]
			(*response)[i].(map[string]interface{})["body"] = []interface{}{map[string]interface{}{
				"payload": payload,
			}}
		}
	case "cloudflare_zero_trust_device_default_profile_local_domain_fallback":
		// wrap the response in 'domains' for tf
		for i := 0; i < resourceCount; i++ {
			do := make(map[string]interface{})
			do["domains"] = []interface{}{(*response)[i]}
			(*response)[i] = do
		}
	case "cloudflare_zero_trust_dex_test":
		denestResponses(response, resourceCount, "dex_tests")
	case "cloudflare_zero_trust_gateway_settings":
		for i := 0; i < resourceCount; i++ {
			settings, ok := (*response)[i].(map[string]interface{})["settings"]
			if !ok {
				return
			}
			customCert, ok := settings.(map[string]interface{})["custom_certificate"]
			if ok {
				delete(customCert.(map[string]interface{}), "binding_status")
				delete(customCert.(map[string]interface{}), "expires_on")
				delete(customCert.(map[string]interface{}), "updated_at")
			}
			blockPage, ok := settings.(map[string]interface{})["block_page"]
			if ok {
				mode := blockPage.(map[string]interface{})["mode"]
				if mode.(string) == "" {
					delete(blockPage.(map[string]interface{}), "mode")
				}
			}
		}
	case "cloudflare_page_rule":
		for i := 0; i < resourceCount; i++ {
			(*response)[i].(map[string]interface{})["target"] = (*response)[i].(map[string]interface{})["targets"].([]interface{})[0].(map[string]interface{})["constraint"].(map[string]interface{})["value"]
			(*response)[i].(map[string]interface{})["actions"] = flattenAttrMap((*response)[i].(map[string]interface{})["actions"].([]interface{}))

			// Have to remap the cache_ttl_by_status to conform to Terraform's more human-friendly structure.
			if cache, ok := (*response)[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_ttl_by_status"].(map[string]interface{}); ok {
				cacheTtlByStatus := []map[string]interface{}{}

				for codes, ttl := range cache {
					if ttl == "no-cache" {
						ttl = 0
					} else if ttl == "no-store" {
						ttl = -1
					}
					elem := map[string]interface{}{
						"codes": codes,
						"ttl":   ttl,
					}

					cacheTtlByStatus = append(cacheTtlByStatus, elem)
				}

				sort.SliceStable(cacheTtlByStatus, func(i int, j int) bool {
					return cacheTtlByStatus[i]["codes"].(string) < cacheTtlByStatus[j]["codes"].(string)
				})

				(*response)[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_ttl_by_status"] = cacheTtlByStatus
			}

			// Remap cache_key_fields.query_string.include & .exclude wildcards (not in an array) to the appropriate "ignore" field value in Terraform.
			if c, ok := (*response)[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_key_fields"].(map[string]interface{}); ok {
				if s, sok := c["query_string"].(map[string]interface{})["include"].(string); sok && s == "*" {
					(*response)[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_key_fields"].(map[string]interface{})["query_string"].(map[string]interface{})["include"] = nil
					(*response)[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_key_fields"].(map[string]interface{})["query_string"].(map[string]interface{})["ignore"] = false
				}
				if s, sok := c["query_string"].(map[string]interface{})["exclude"].(string); sok && s == "*" {
					(*response)[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_key_fields"].(map[string]interface{})["query_string"].(map[string]interface{})["exclude"] = nil
					(*response)[i].(map[string]interface{})["actions"].(map[string]interface{})["cache_key_fields"].(map[string]interface{})["query_string"].(map[string]interface{})["ignore"] = true
				}
			}
		}
	case "cloudflare_zero_trust_access_short_lived_certificate":
		remapProperty(response, resourceCount, "id", "app_id")
	case "cloudflare_zone_setting":
		remapProperty(response, resourceCount, "id", "setting_id")
	case "cloudflare_hostname_tls_setting":
		addAttributeKeyValue(response, resourceCount, "setting_id", pathParam)
	case "cloudflare_registrar_domain":
		remapProperty(response, resourceCount, "name", "domain_name")
	case "cloudflare_r2_managed_domain":
		addAttributeKeyValue(response, resourceCount, "bucket_name", pathParam)
	case "cloudflare_r2_custom_domain":
		finalResponse := make([]interface{}, 0)
		r := *response
		for i := 0; i < resourceCount; i++ {
			domains := r[i].(map[string]interface{})["domains"]
			bucketObjects := make([]interface{}, len(domains.([]interface{})))
			for j := range domains.([]interface{}) {
				b := domains.([]interface{})[j]
				b.(map[string]interface{})["bucket_name"] = pathParam
				b.(map[string]interface{})["zone_id"] = b.(map[string]interface{})["zoneId"]
				bucketObjects[j] = b
			}
			finalResponse = append(finalResponse, bucketObjects...)
		}
		*response = make([]interface{}, len(finalResponse))
		for i := range finalResponse {
			(*response)[i] = finalResponse[i]
		}
	case "cloudflare_pages_domain":
		addAttributeKeyValue(response, resourceCount, "project_name", pathParam)
	case "cloudflare_list_item":
		remapProperty(response, resourceCount, "id", "list_id")
	case "cloudflare_api_shield_schema":
		remapProperty(response, resourceCount, "source", "file")
	case "cloudflare_api_shield_discovery_operation":
		remapProperty(response, resourceCount, "id", "operation_id")
	case "cloudflare_zero_trust_dlp_predefined_profile":
		addAttributeKeyValue(response, resourceCount, "profile_id", pathParam)
	case "cloudflare_zero_trust_access_identity_provider":
		for i := 0; i < resourceCount; i++ {
			cfg, ok := (*response)[i].(map[string]interface{})["config"]
			if ok {
				delete(cfg.(map[string]interface{}), "redirect_url")
			}
			scimCFG, ok := (*response)[i].(map[string]interface{})["scim_config"]
			if ok {
				delete(scimCFG.(map[string]interface{}), "scim_base_url")
			}
		}
	case "cloudflare_zero_trust_access_custom_page":
		// fetch each object one by one to get 'custom_html' field.
		endpointFMT := resourceToEndpoint[resourceType]["get"]
		placeholderReplacer := strings.NewReplacer("{account_id}", accountID)
		endpointFMT = placeholderReplacer.Replace(endpointFMT)
		for i := 0; i < resourceCount; i++ {
			uid, ok := (*response)[i].(map[string]interface{})["uid"]
			if !ok {
				continue
			}
			endpoint := strings.Replace(endpointFMT, "{custom_page_id}", uid.(string), 1)
			result := new(http.Response)
			err := api.Get(context.Background(), endpoint, nil, &result)
			if err != nil {
				var apierr *cloudflare.Error
				if errors.As(err, &apierr) {
					if apierr.StatusCode == http.StatusNotFound {
						log.WithFields(logrus.Fields{
							"resource": resourceType,
							"endpoint": endpoint,
						}).Debug("no resources found")
					}
				}
				log.Fatalf("failed to fetch API endpoint: %s", err)
			}
			body, err := io.ReadAll(result.Body)
			if err != nil {
				log.Fatalln(err)
			}
			value := gjson.Get(string(body), "result")
			if value.Type == gjson.Null {
				log.WithFields(logrus.Fields{
					"resource": resourceType,
					"endpoint": endpoint,
				}).Debug("no result found")
				continue
			}
			customHTML := gjson.Get(value.Raw, "custom_html")
			if value.Type == gjson.Null {
				continue
			}
			(*response)[i].(map[string]interface{})["custom_html"] = customHTML.String()
		}
	case "cloudflare_web_analytics_rule":
		finalResponse := make([]interface{}, 0)
		r := *response
		for i := 0; i < resourceCount; i++ {
			rules := r[i].(map[string]interface{})["rules"]
			ruleObjects := make([]interface{}, len(rules.([]interface{})))
			for j := range rules.([]interface{}) {
				b := rules.([]interface{})[j]
				b.(map[string]interface{})["ruleset_id"] = pathParam
				ruleObjects[j] = b
			}
			finalResponse = append(finalResponse, ruleObjects...)
		}
		*response = make([]interface{}, len(finalResponse))
		for i := range finalResponse {
			(*response)[i] = finalResponse[i]
		}
	case "cloudflare_waiting_room_event":
		addAttributeKeyValue(response, resourceCount, "waiting_room_id", pathParam)
	case "cloudflare_waiting_room_rules":
		*response = []interface{}{
			map[string]interface{}{
				"waiting_room_id": pathParam,
				"rules":           *response,
			},
		}
	case "cloudflare_keyless_certificate":
		addAttributeKeyValue(response, resourceCount, "certificate", "-----INSERT CERTIFICATE-----")
	case "cloudflare_stream_watermark":
		addAttributeKeyValue(response, resourceCount, "file", `REPLACE with filebase64("path-to-file")`)
	case "cloudflare_authenticated_origin_pulls_certificate":
		addAttributeKeyValue(response, resourceCount, "private_key", "-----INSERT PRIVATE KEY-----")
	case "cloudflare_zero_trust_access_mtls_certificate":
		addAttributeKeyValue(response, resourceCount, "certificate", "-----INSERT CERTIFICATE-----")
	case "cloudflare_zero_trust_access_mtls_hostname_settings":
		*response = []interface{}{
			map[string]interface{}{
				"settings": *response,
			},
		}
	case "cloudflare_workers_script_subdomain":
		addAttributeKeyValue(response, resourceCount, "script_name", pathParam)
	case "cloudflare_workers_deployment":
		finalResponse := make([]interface{}, 0)
		r := *response
		for i := 0; i < resourceCount; i++ {
			deployments := r[i].(map[string]interface{})["deployments"]
			deploymentObjects := make([]interface{}, len(deployments.([]interface{})))
			for j := range deployments.([]interface{}) {
				d := deployments.([]interface{})[j]
				d.(map[string]interface{})["script_name"] = pathParam
				deploymentObjects[j] = d
			}
			finalResponse = append(finalResponse, deploymentObjects...)
		}
		*response = make([]interface{}, len(finalResponse))
		for i := range finalResponse {
			(*response)[i] = finalResponse[i]
		}
	case "cloudflare_workers_cron_trigger":
		for i := 0; i < resourceCount; i++ {
			(*response)[i].(map[string]interface{})["script_name"] = pathParam
			schedules, ok := (*response)[i].(map[string]interface{})["schedules"]
			if !ok {
				continue
			}
			for j := range schedules.([]interface{}) {
				delete(schedules.([]interface{})[j].(map[string]interface{}), "created_on")
				delete(schedules.([]interface{})[j].(map[string]interface{}), "modified_on")
			}
		}
	case "cloudflare_authenticated_origin_pulls":
		for i := 0; i < resourceCount; i++ {
			hName := (*response)[i].(map[string]interface{})["hostname"]
			cID := (*response)[i].(map[string]interface{})["cert_id"]
			enabled := (*response)[i].(map[string]interface{})["enabled"]
			(*response)[i].(map[string]interface{})["config"] = []interface{}{
				map[string]interface{}{
					"hostname": hName,
					"cert_id":  cID,
					"enabled":  enabled,
				},
			}
		}
	case "cloudflare_magic_wan_static_route":
		denestResponses(response, resourceCount, "routes")
	case "cloudflare_ruleset":
		ruleHeaders := map[string][]map[string]interface{}{}
		for i, ruleset := range *response {
			if ruleset.(map[string]interface{})["rules"] != nil {
				for j, rule := range ruleset.(map[string]interface{})["rules"].([]interface{}) {
					ID := rule.(map[string]interface{})["id"]
					if ID != nil {
						headers, exists := ruleHeaders[ID.(string)]
						if exists {
							(*response)[i].(map[string]interface{})["rules"].([]interface{})[j].(map[string]interface{})["action_parameters"].(map[string]interface{})["headers"] = headers
						}
					}
				}
			}
		}

		// log custom fields specific transformation fields
		logCustomFieldsTransform := []string{"cookie_fields", "request_fields", "response_fields"}

		for i := 0; i < resourceCount; i++ {
			rules := (*response)[i].(map[string]interface{})["rules"]
			if rules != nil {
				for ruleCounter := range rules.([]interface{}) {
					// should the `ref` be the default `id`, don't output it
					// as we don't need to track a computed default.
					id := rules.([]interface{})[ruleCounter].(map[string]interface{})["id"]
					ref := rules.([]interface{})[ruleCounter].(map[string]interface{})["ref"]
					if id == ref {
						rules.([]interface{})[ruleCounter].(map[string]interface{})["ref"] = nil
					}

					actionParams := rules.([]interface{})[ruleCounter].(map[string]interface{})["action_parameters"]
					if actionParams != nil {
						// check for log custom fields that need to be transformed
						for _, logCustomFields := range logCustomFieldsTransform {
							// check if the field exists and make sure it has at least one element
							if actionParams.(map[string]interface{})[logCustomFields] != nil && len(actionParams.(map[string]interface{})[logCustomFields].([]interface{})) > 0 {
								// Create a new list to store the data in.
								var newLogCustomFields []interface{}
								// iterate over each of the keys and add them to a generic list
								for logCustomFieldsCounter := range actionParams.(map[string]interface{})[logCustomFields].([]interface{}) {
									newLogCustomFields = append(newLogCustomFields, map[string]interface{}{"name": actionParams.(map[string]interface{})[logCustomFields].([]interface{})[logCustomFieldsCounter].(map[string]interface{})["name"]})
								}
								actionParams.(map[string]interface{})[logCustomFields] = newLogCustomFields
							}
						}

						// check if our ruleset is of action 'skip'
						if rules.([]interface{})[ruleCounter].(map[string]interface{})["action"] == "skip" && (*response)[i].(map[string]interface{})["phase"] != "http_request_firewall_managed" {
							for rule := range actionParams.(map[string]interface{}) {
								// "rules" is the only map[string][]string we need to remap. The others are all []string and are handled naturally.
								if rule == "rules" {
									for key, value := range actionParams.(map[string]interface{})[rule].(map[string]interface{}) {
										var rulesList []string
										for _, val := range value.([]interface{}) {
											rulesList = append(rulesList, val.(string))
										}
										actionParams.(map[string]interface{})[rule].(map[string]interface{})[key] = strings.Join(rulesList, ",")
									}
								}
							}
						}

						// Cache Rules transformation
						if (*response)[i].(map[string]interface{})["phase"] == "http_request_cache_settings" {
							if ck, ok := rules.([]interface{})[ruleCounter].(map[string]interface{})["action_parameters"].(map[string]interface{})["cache_key"]; ok {
								if c, cok := ck.(map[string]interface{})["custom_key"]; cok {
									if qs, qok := c.(map[string]interface{})["query_string"]; qok {
										if s, sok := qs.(map[string]interface{})["include"]; sok && s == "*" {
											rules.([]interface{})[ruleCounter].(map[string]interface{})["action_parameters"].(map[string]interface{})["cache_key"].(map[string]interface{})["custom_key"].(map[string]interface{})["query_string"].(map[string]interface{})["include"] = map[string]interface{}{"list": []string{"*"}}
										}
										if s, sok := qs.(map[string]interface{})["exclude"]; sok && s == "*" {
											rules.([]interface{})[ruleCounter].(map[string]interface{})["action_parameters"].(map[string]interface{})["cache_key"].(map[string]interface{})["custom_key"].(map[string]interface{})["query_string"].(map[string]interface{})["exclude"] = map[string]interface{}{"list": []string{"*"}}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func unMarshallJSONStructData(modifiedJSONString string) ([]interface{}, error) {
	var data interface{}
	err := json.Unmarshal([]byte(modifiedJSONString), &data)
	if err != nil {
		return nil, err
	}
	if dataSlice, ok := data.([]interface{}); ok {
		return dataSlice, nil
	}
	return []interface{}{data}, nil
}

func getAPIResponse(result *http.Response, pathParams []string, endpoints ...string) ([]interface{}, error) {
	var allResults []interface{}

	for i, baseEndpoint := range endpoints {
		page := 1
		totalPages := 1
		param := ""
		if len(pathParams) > 0 {
			param = pathParams[i]
		}

		for {
			var endpoint string
			// no page param for first request
			if page == 1 {
				endpoint = baseEndpoint
			} else {
				sep := "?"
				if strings.Contains(baseEndpoint, "?") {
					sep = "&"
				}
				endpoint = fmt.Sprintf("%s%spage=%d", baseEndpoint, sep, page)
			}

			err := api.Get(context.Background(), endpoint, nil, &result)
			if err != nil {
				var apierr *cloudflare.Error
				if errors.As(err, &apierr) && apierr.StatusCode == http.StatusNotFound {
					log.WithFields(logrus.Fields{
						"resource": resourceType,
						"endpoint": endpoint,
					}).Debug("no resources found")
					return nil, err
				}
				log.Fatalf("failed to fetch API endpoint: %s", err)
			}

			body, err := io.ReadAll(result.Body)
			if err != nil {
				log.Fatalln(err)
			}

			resultVal := gjson.Get(string(body), "result")
			if resultVal.Type == gjson.Null {
				log.WithFields(logrus.Fields{
					"resource": resourceType,
					"endpoint": endpoint,
				}).Debug("no result found")
				return nil, errors.New("no result found")
			}

			modifiedJSON := modifyResponsePayload(resourceType, resultVal)
			jsonStructData, err := unMarshallJSONStructData(modifiedJSON)
			if err != nil {
				log.Fatalf("failed to unmarshal result: %s", err)
			}

			processCustomCasesV5(&jsonStructData, resourceType, param)
			allResults = append(allResults, jsonStructData...)

			if page == 1 {
				totalPagesVal := gjson.Get(string(body), "result_info.total_pages")
				if totalPagesVal.Exists() {
					totalPages = int(totalPagesVal.Int())
				}
			}

			if page >= totalPages {
				break
			}
			page++
		}
	}
	return allResults, nil
}

func isSupportedPathParam(resources []string, rType string) bool {
	_, ok := settingsMap[rType]
	if !ok {
		return false
	}
	return slices.Contains(resources, rType)
}

func replacePathParams(params []string, endpoint string, rType string) []string {
	endpoints := make([]string, 0)
	var placeholder string
	switch rType {
	case "cloudflare_zone_setting", "cloudflare_hostname_tls_setting":
		placeholder = "{setting_id}"
	case "cloudflare_waiting_room_event":
		placeholder = "{waiting_room_id}"
	case "cloudflare_r2_managed_domain", "cloudflare_r2_custom_domain":
		placeholder = "{bucket_name}"
	case "cloudflare_pages_domain":
		placeholder = "{project_name}"
	case "cloudflare_list_item":
		placeholder = "{list_id}"
	case "cloudflare_zero_trust_dlp_predefined_profile":
		placeholder = "{profile_id}"
	case "cloudflare_web_analytics_rule":
		placeholder = "{ruleset_id}"
	case "cloudflare_waiting_room_rules":
		placeholder = "{waiting_room_id}"
	case "cloudflare_zero_trust_tunnel_cloudflared_config":
		placeholder = "{tunnel_id}"
	case "cloudflare_workers_script_subdomain", "cloudflare_workers_deployment", "cloudflare_workers_cron_trigger":
		placeholder = "{script_name}"
	case "cloudflare_authenticated_origin_pulls":
		placeholder = "{hostname}"
	case "cloudflare_queue_consumer":
		placeholder = "{queue_id}"
	case "cloudflare_api_shield_operation_schema_validation_settings":
		placeholder = "{operation_id}"
	case "cloudflare_observatory_scheduled_test":
		for _, id := range params {
			endpoints = append(endpoints, strings.Clone(strings.NewReplacer("{url}", url.QueryEscape(id)).Replace(endpoint)))
		}
		return endpoints
	case "cloudflare_zero_trust_dlp_custom_profile":
		placeholder = "{profile_id}"
	default:
		return endpoints
	}
	for _, id := range params {
		endpoints = append(endpoints, strings.Clone(strings.NewReplacer(placeholder, id).Replace(endpoint)))
	}
	return endpoints
}

func addAttributeKeyValue(response *[]interface{}, resourceCount int, key string, value string) {
	for i := 0; i < resourceCount; i++ {
		(*response)[i].(map[string]interface{})[key] = value
	}
}

func remapProperty(response *[]interface{}, resourceCount int, responseProperty string, remappedProperty string) {
	for i := 0; i < resourceCount; i++ {
		prop, ok := (*response)[i].(map[string]interface{})[responseProperty]
		if !ok {
			continue
		}
		(*response)[i].(map[string]interface{})[remappedProperty] = prop
	}
}

func denestResponses(response *[]interface{}, resourceCount int, nestedAttributeName string) {
	finalResponse := make([]interface{}, 0)
	r := *response
	for i := 0; i < resourceCount; i++ {
		nestedObjects := r[i].(map[string]interface{})[nestedAttributeName]
		objects := make([]interface{}, len(nestedObjects.([]interface{})))
		for j := range nestedObjects.([]interface{}) {
			o := nestedObjects.([]interface{})[j]
			objects[j] = o
		}
		finalResponse = append(finalResponse, objects...)
	}
	*response = make([]interface{}, len(finalResponse))
	for i := range finalResponse {
		(*response)[i] = finalResponse[i]
	}
}

func processRulesetV5() {

}
