#!/bin/bash -e

. private/application-id.sh

cd cmd/wiki-skill-server
export cert=../../private/certificate.pem
export key=../../private/private-key.pem
fresh
