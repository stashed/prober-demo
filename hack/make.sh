#!/bin/bash
set -xeou pipefail

GOPATH=$(go env GOPATH)
REPO_ROOT=$GOPATH/src/stash.appscode.dev/prober-demo
REGISTRY=appscodeci
BINARY_NAME=prober-demo

# build binary
pushd $REPO_ROOT
go build -ldflags "-linkmode external -extldflags -static" -o $BINARY_NAME main.go
chmod +x $BINARY_NAME
docker build -t $REGISTRY/$BINARY_NAME .

popd

docker push $REGISTRY/$BINARY_NAME