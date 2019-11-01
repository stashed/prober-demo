#!/bin/bash
set -xeou pipefail

GOPATH=$(go env GOPATH)
REPO_ROOT=$GOPATH/src/github.com/tamalsaha/prober-demo
REGISTRY=emruzhossain
BINARY_NAME=prober-demo

# build binary
pushd $REPO_ROOT
go build -ldflags "-linkmode external -extldflags -static" -o $BINARY_NAME main.go
chmod +x $BINARY_NAME
docker build -t $REGISTRY/$BINARY_NAME .

popd

# load image into kind cluster
kind load docker-image $REGISTRY/$BINARY_NAME