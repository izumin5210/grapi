#!/usr/bin/env bash

set -eu
set -o pipefail

cd "$(dirname $0)/../.."

DIR=_tests/integration
IMAGE_NAME=grapi-integration-test:go$GO_VERSION

TARGET_REVISION=${TRAVIS_COMMIT:-$(git symbolic-ref --short HEAD)}

docker build -t $IMAGE_NAME --build-arg GO_VERSION=$GO_VERSION -f ./$DIR/Dockerfile .
docker run \
  -v $(pwd)/$DIR:/go/src/e2e \
  --env TARGET_REVISION=$TARGET_REVISION \
  $IMAGE_NAME \
  sh -c 'go test -v -timeout 2m'
