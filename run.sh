#!/bin/bash
echo "This is an example - customise to suit your own use case"
echo "It is assumed you are already running an instance of session host"
echo "To start the rfswitch's USB connection to session host type"
echo "./arduino/test.sh ttyUSB0 &"
echo "to starting the calibration service:"
echo "cd ./py/ && nohup python client.py &"
echo "check out the error messages if not working with websocat ws://localhost:8888/ws/data -  "
echo "Setting environment variables"
export VNA_DESTINATION=ws://localhost:8888/ws/data
export VNA_RFSWITCH=ws://localhost:8888/ws/rfswitch
export VNA_CALIBRATION=ws://localhost:8888/ws/calibration
export VNA_DEVELOPMENT=true
./cmd/vna/vna unlock 
./cmd/vna/vna stream 
echo "You can manually test the system using: websocat ws://localhost:8888/ws/data -"
echo "Or run ./py/validate.py to perform calibration, measure DUT and plot results"

