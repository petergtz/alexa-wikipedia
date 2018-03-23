#!/bin/bash -e

. private/dev-application-id.sh

export cert=private/certificate.pem
export key=private/private-key.pem
fresh
