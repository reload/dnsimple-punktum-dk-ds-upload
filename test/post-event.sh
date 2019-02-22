#!/bin/bash

TOKEN=$(sed -n -e '/^TOKEN:/ s/.*: *//p' "$(dirname "$0")/../env.yaml")

URL="https://europe-west1-reload-internal-alpha.cloudfunctions.net/dnsimple-dk-hostmaster-ds-upload?token={$TOKEN}"

curl -D - -X POST -H "Content-Type: application/json" -d "@$(dirname "$0")/dnssec.rotation_complete.json" "$URL"
