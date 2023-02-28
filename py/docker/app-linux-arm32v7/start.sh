#!/bin/bash
app="calibration-twport"
docker build -t ${app} .
docker run ${app} 

