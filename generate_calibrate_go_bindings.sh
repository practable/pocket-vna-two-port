#!/bin/bash
 protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=./pkg/pb --go-grpc_opt=paths=source_relative \
    calibrate.proto
