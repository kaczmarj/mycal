#!/usr/bin/env bash

GOARCH="amd64"
for GOOS in darwin linux windows;
do
    output="mycal-${GOOS}-${GOARCH}"
    if [ $GOOS = "windows" ]; then
        output="${output}.exe"
    fi
    echo "Compiling $output ..."
    env GOOS="$GOOS" GOARCH="$GOARCH" go build -o "$output" .
done
