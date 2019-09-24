#!/bin/bash -e

cd $(dirname $0)/..

for region in 'us-east-1' 'eu-west-1' 'ap-northeast-1'; do
    output=$(aws --region $region lambda publish-version --function-name AlexaWikipedia)
    echo $output
    version=$(echo $output | jq -r .Version)
    sed -E -i "s/(arn:aws:lambda:$region:512841817041:function:AlexaWikipedia:)([0-9]+)/\1$version/g" skill.json
done

ask diff --target skill

while true; do
    read -p "Deploy changes? " yn
    case $yn in
        [Yy]* ) break;;
        [Nn]* ) exit 1;;
        * ) echo "Please answer yes or no.";;
    esac
done

ask deploy --force --target skill

### ask api submit