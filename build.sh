#!/bin/zsh

# Build
pushd web || exit
npm run build
popd || exit

go build -ldflags "-s -w"
echo build completed