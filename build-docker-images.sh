#!/usr/bin/env bash

# Get version
VERSION=$(grep 'version string = ' engine/version.go | sed 's:^[^"]*"\([^"]*\)".*:\1:')

# Build Go 'rules-engine'
docker build --rm -t rules-engine .
docker tag rules-engine mainflux-labs/rules-engine:${VERSION}
docker tag rules-engine mainflux-labs/rules-engine:latest

# Build Python 'parser'
cd ./parser/
docker build --rm -t rules-engine-parser .
docker tag rules-engine-parser mainflux-labs/rules-engine-parser:${VERSION}
docker tag rules-engine-parser mainflux-labs/rules-engine-parser:latest
