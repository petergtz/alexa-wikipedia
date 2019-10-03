#!/bin/bash -ex

cd $(dirname $0)/..

for region in 'us-east-1' 'eu-west-1' 'ap-northeast-1'; do
    output=$(aws --region $region lambda publish-version --function-name AlexaWikipedia)
    echo $output
    version=$(echo $output | jq -r .Version)
    
    aws --region $region lambda update-alias \
      --function-name AlexaWikipedia \
      --name prod \
      --function-version $version

    aws --region $region lambda add-permission \
      --function-name AlexaWikipedia:prod \
      --action lambda:invokeFunction \
      --principal alexa-appkit.amazon.com  \
      --statement-id $(date +%s) \
      --event-source-token amzn1.ask.skill.101a1dcb-ce70-4fd6-ae01-6a8803f727ff
done
