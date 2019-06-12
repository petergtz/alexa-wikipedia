#!/bin/bash -e

. private/dev-application-id.sh
. private/s3-credentials.sh

export PORT=4443
export SKILL_USE_TLS=true
export CERT=private/certificate.pem
export KEY=private/private-key.pem
export TABLE_NAME_OVERRIDE=TestAlexaWikipediaRequests
fresh
