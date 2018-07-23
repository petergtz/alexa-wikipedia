#!/bin/bash -ex

cd $(dirname $0)/..

if [[ -n $(git status -s) ]]; then
    echo "Working directory not clean"
    exit 1
fi

ginkgo -r

export SHA=$(git rev-parse --short HEAD)

cf push alexa-wikipedia-$SHA
