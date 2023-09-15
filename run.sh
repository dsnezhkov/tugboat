#!/bin/bash

export GCC=x86_64-w64-mingw32-gcc
CC=$GCC CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build  --ldflags="-s -w" -o ./bin/tugboat.exe . && \
 arch -x86_64  /opt/homebrew/bin/wine64 ./bin/tugboat.exe  2>&1
