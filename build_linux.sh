#!/usr/bin/env bash
source env.sh
#
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

export SKIP_TEST=true

##
export EXEFILE=setup_linux

./build.sh