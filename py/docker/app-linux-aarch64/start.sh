#!/bin/bash
app="calibration"
docker build -t ${app} .
docker run ${app} 

