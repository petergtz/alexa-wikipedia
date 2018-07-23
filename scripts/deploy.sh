#!/bin/bash -ex

cd $(dirname $0)/..

if [[ -n $(git status -s) ]]; then
    echo "Working directory not clean"
    exit 1
fi

ginkgo -r

export SHA=$(git rev-parse --short HEAD)

cf push alexa-wikipedia-$SHA

export OLD_RELEASES=$(cf apps | grep alexa-wikipedia | grep -v $SHA | cut -f 1 -d ' ')

for release in $OLD_RELEASES; do
    cf delete -f $release
done
