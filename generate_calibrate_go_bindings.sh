#!/bin/bash
#needs recent version of grpcio to match container image, e.g. conda activate py3.12
protoc --go_out=./pkg/pb --go_opt=paths=source_relative \
   --go-grpc_out=./pkg/pb --go-grpc_opt=paths=source_relative \
   calibrate.proto
