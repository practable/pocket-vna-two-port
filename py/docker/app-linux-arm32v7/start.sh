#!/bin/bash
app="calibration-twoport"
docker build -t ${app} .
docker run ${app} 

