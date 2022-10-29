#!/bin/zsh

# Build
go build

# Run
./trade-ripper -q home.docker.lan:9009 -c crypto
