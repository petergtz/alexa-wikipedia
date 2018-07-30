#!/bin/bash -ex

cd $(dirname $0)/..

if [[ -n $(git status -s) ]]; then
    echo "Working directory not clean"
    exit 1
fi

ginkgo -r

export SHA=$(git rev-parse --short HEAD)
export APP_NAME=alexa-wikipedia-$SHA

cf push --no-start -b https://github.com/cloudfoundry/go-buildpack.git $APP_NAME
cf set-env $APP_NAME APPLICATION_ID $(lpass show Personal/Alexa-Wikipedia-Skill-Application-ID --password)
cf restart $APP_NAME

export OLD_RELEASES=$(cf apps | grep alexa-wikipedia | grep -v $SHA | cut -f 1 -d ' ')

for release in $OLD_RELEASES; do
    cf delete -f $release
done
