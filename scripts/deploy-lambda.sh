#!/bin/bash -ex

cd $(dirname $0)/..

go build -o main
rm -f main.zip
zip alexa-wikipedia.zip main
aws s3 cp alexa-wikipedia.zip s3://alexa-golang-skills/alexa-wikipedia.zip
aws lambda update-function-code --function-name AlexaWikipedia --s3-bucket alexa-golang-skills --s3-key alexa-wikipedia.zip
rm -f alexa-wikipedia.zip
rm -f main