#!/bin/bash

TOKEN=$(yq read "$(dirname "$0")/../env.yaml" TOKEN)

URL=$(gcloud functions describe dnsimple-dk-hostmaster-ds-upload --region=europe-west1 | yq read - httpsTrigger.url)

curl -D - -X POST -H "Content-Type: application/json" -d "@$(dirname "$0")/dnssec.rotation_complete.json" "${URL}?token=${TOKEN}"
