#!/bin/bash

set -eou pipefail

tmpd=$(mktemp -d)
go build -o $tmpd .
(
    cd $tmpd
    echo "provider ${PROVIDER_NAME} {}" > main.tf

    docker run \
        --interactive \
        --rm \
        --tty \
        --volume "${tmpd}:/out" \
        --workdir /out \
        hashicorp/terraform \
        init 2>&1 >/dev/null

    docker run \
        --interactive \
        --rm \
        --tty \
        --volume "${tmpd}:/out" \
        --workdir /out \
        hashicorp/terraform \
        providers schema -json > schema.json

    rm main.tf
    rm terraform-provider-${PROVIDER_NAME}
)
echo $tmpd
