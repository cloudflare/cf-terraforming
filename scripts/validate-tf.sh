#!/usr/bin/env bash

set -e
export TF_IN_AUTOMATION=1

rootdir=$(pwd)
for dir in testdata/terraform/v4/*; do
    if [ $dir == "testdata/terraform/v4/cloudflare_load_balancer" ]; then
        continue
    fi

    echo "==> $dir (test.tf)"
    cd $dir
    terraform init -backend=false -no-color
    terraform validate -no-color
    cd $rootdir
done
