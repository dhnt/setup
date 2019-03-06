#!/usr/bin/env bash
source env.sh
#
export GOOS=darwin
export GOARCH=amd64
export CGO_ENABLED=0

export SKIP_TEST=false

#
export EXEFILE=setup_darwin

./build.sh