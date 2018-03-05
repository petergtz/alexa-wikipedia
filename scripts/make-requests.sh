#!/bin/bash -e
. private/endpoint.sh

curl -X POST -d "@private/IntentRequest.json" $ENDPOINT --cacert private/certificate.pem
curl -X POST -d "@private/LaunchRequest.json" $ENDPOINT --cacert private/certificate.pem