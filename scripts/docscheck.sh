#!/bin/bash

set -eou pipefail

echo "==> Extracting provider json schema..."
tmpd=$(mktemp -d)
go build -o "${tmpd}" .
(
    cd "${tmpd}"
    echo "provider ${PROVIDER_NAME} {}" > main.tf

    docker run \
        --interactive \
        --rm \
        --tty \
        --volume "${tmpd}:/out" \
        --workdir /out \
        hashicorp/terraform \
        init >/dev/null 2>&1

    docker run \
        --interactive \
        --rm \
        --tty \
        --volume "${tmpd}:/out" \
        --workdir /out \
        hashicorp/terraform \
        providers schema -json >schema.json

    sudo rm -rf .terraform
    rm main.tf "terraform-provider-${PROVIDER_NAME}"
)

echo "==> Checking docs with tfproviderdocs..."

docker run \
    --interactive \
    --rm \
    --tty \
    --volume "${PWD}/website:/terraform-provider-megaport/website" \
    --volume "${tmpd}:/provider-schema" \
    --workdir /terraform-provider-megaport \
    bflad/tfproviderdocs \
    check \
    -allowed-resource-subcategories=resources,datasources \
    -providers-schema-json=/provider-schema/schema.json \
    -require-resource-subcategory

rm -rf "${tmpd}"
