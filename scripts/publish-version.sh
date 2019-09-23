#!/bin/bash -ex

cd $(dirname $0)/..

aws lambda publish-version --function-name AlexaWikipedia
