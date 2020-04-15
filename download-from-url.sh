#!/bin/sh
URL="$1"

GCS_PREFIX="https://prow.svc.ci.openshift.org/view/gcs/origin-ci-test/"
if [ "${URL#$GCS_PREFIX}" != "$URL" ]; then
    DIR="./output/$(echo "${URL#$GCS_PREFIX}" | sed 's/[^0-9A-Za-z._-]/__/g')/"
    PREFIX="${URL#$GCS_PREFIX}/"
    mkdir -p "$DIR"
    ./list-bucket -output-dir="$DIR" -prefix="$PREFIX"
fi
