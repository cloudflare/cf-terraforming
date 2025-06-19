#!/usr/bin/env bash

set -e
export TF_IN_AUTOMATION=1

SKIP_DIRS=(
"testdata/terraform/v5/cloudflare_api_token"
"testdata/terraform/v5/cloudflare_access_application"
"testdata/terraform/v5/cloudflare_access_identity_provider"
"testdata/terraform/v5/cloudflare_access_identity_provider"
"testdata/terraform/v5/cloudflare_access_rule"
"testdata/terraform/v5/cloudflare_argo"
"testdata/terraform/v5/cloudflare_custom_pages"
"testdata/terraform/v5/cloudflare_cloud_connector_rules"
"testdata/terraform/v5/cloudflare_cloudforce_one_request"
"testdata/terraform/v5/cloudflare_cloudforce_one_request_asset"
"testdata/terraform/v5/cloudflare_cloudforce_one_request_message"
"testdata/terraform/v5/cloudflare_cloudforce_one_request_priority"
"testdata/terraform/v5/cloudflare_custom_ssl"
"testdata/terraform/v5/cloudflare_magic_network_monitoring_configuration"
"testdata/terraform/v5/cloudflare_magic_network_monitoring_rule"
"testdata/terraform/v5/cloudflare_magic_transit_connector"
"testdata/terraform/v5/cloudflare_magic_transit_site"
"testdata/terraform/v5/cloudflare_magic_transit_site_acl"
"testdata/terraform/v5/cloudflare_magic_transit_site_lan"
"testdata/terraform/v5/cloudflare_magic_transit_site_wan"
"testdata/terraform/v5/cloudflare_magic_wan_gre_tunnel"
"testdata/terraform/v5/cloudflare_magic_wan_ipsec_tunnel"
"testdata/terraform/v5/cloudflare_firewall_rule"
"testdata/terraform/v5/cloudflare_image"
"testdata/terraform/v5/cloudflare_image_variant"
"testdata/terraform/v5/cloudflare_workers_script"
"testdata/terraform/v5/cloudflare_workers_secret"
"testdata/terraform/v5/cloudflare_zone_subscription"
)
rootdir=$(pwd)

for dir in testdata/terraform/v5/*; do
    skip=false
    for skipdir in "${SKIP_DIRS[@]}"; do
        if [[ "$dir" == "$skipdir" ]]; then
            skip=true
            break
        fi
    done

    if $skip; then
        echo "==> Skipping $dir"
        continue
    fi
    echo "==> $dir (test.tf)"
    cd $dir
    terraform init -backend=false -no-color
    terraform validate -no-color
    cd $rootdir
done
