#!/bin/bash -ex

cd $(dirname $0)/..

pushd cmd/skill
    GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -o bootstrap


    rm -f alexa-wikipedia.zip
    zip alexa-wikipedia.zip bootstrap

    aws s3 cp alexa-wikipedia.zip s3://alexa-golang-skills/alexa-wikipedia.zip
    aws s3 cp s3://alexa-golang-skills/alexa-wikipedia.zip s3://alexa-golang-skills-eu-west-1/alexa-wikipedia.zip &
    aws s3 cp s3://alexa-golang-skills/alexa-wikipedia.zip s3://alexa-golang-skills-ap-northeast-1/alexa-wikipedia.zip &
    wait
    rm -f alexa-wikipedia.zip
    rm -f bootstrap
popd

aws --region us-east-1 lambda update-function-code --function-name AlexaWikipedia --s3-bucket alexa-golang-skills --s3-key alexa-wikipedia.zip &
aws --region eu-west-1 lambda update-function-code --function-name AlexaWikipedia --s3-bucket alexa-golang-skills-eu-west-1 --s3-key alexa-wikipedia.zip &
aws --region ap-northeast-1 lambda update-function-code --function-name AlexaWikipedia --s3-bucket alexa-golang-skills-ap-northeast-1 --s3-key alexa-wikipedia.zip &
wait

