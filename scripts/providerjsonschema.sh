#!/bin/bash

set -eou pipefail

tmpd=$(mktemp -d)
go build -o $tmpd .
(
    cd $tmpd
    echo "provider ${PROVIDER_NAME} {}" > main.tf
    terraform init 2>&1 >/dev/null
    terraform providers schema -json > schema.json
    rm main.tf
    rm terraform-provider-${PROVIDER_NAME}
)
echo $tmpd
