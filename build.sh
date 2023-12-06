#!/bin/bash
export GOOS=linux
now=$(date +'%Y-%m-%d_%T') #no spaces to prevent quotation issues in build command
(cd ./cmd/vna; go build -ldflags "-X 'github.com/practable/pocket-vna-two-port/cmd/vna/cmd.Version=`git describe --tags`' -X 'github.com/practable/pocket-vna-two-port/cmd/vna/cmd.BuildTime=$now'"; ./vna version)
