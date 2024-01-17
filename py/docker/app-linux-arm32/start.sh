#!/bin/bash
app="calibration-twoport-grpc"
docker build -t ${app} .
docker run ${app} 

