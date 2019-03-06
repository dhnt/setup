#!/usr/bin/env bash

##
source env.sh
##
[[ $DEBUG ]] && FLAG="-x"

function build() {
    echo "## Cleaning ..."
    go clean $FLAG ./...

    echo "## Formatting ..."
    go fmt $FLAG ./...; if [ $? -ne 0 ]; then
        return 1
    fi
    
    echo "## Vetting ..."
    go vet $FLAG ./...; if [ $? -ne 0 ]; then
        return 1
    fi

    echo "## Testing ..."
    if [ "x${SKIP_TEST}" != "xtrue" ]; then
        go test $FLAG ./...; if [ $? -ne 0 ]; then
            return 1
        fi
    fi

    echo "## Building ..."
    
    go build $FLAG -o $EXEFILE -a -ldflags '-w -extldflags "-static"' ./...; if [ $? -ne 0 ]; then
        return 1
    fi

    echo "## Tidying up modules ..."
    go mod tidy
}

echo "#### Building ..."

build; if [ $? -ne 0 ]; then
    echo "#### Build failure"
    exit 1
fi

echo "#### Build success"

exit 0
