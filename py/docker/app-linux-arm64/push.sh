#!/bin/bash
docker tag calibration-twoport:latest practable/calibration-twoport-grpc:arm64v8-3.12-0.2
docker push practable/calibration-twoport-grpc:arm64v8-3.12-0.2
