#!/bin/bash

# we assume jump is already installed
# we assume that relay access and token files are already present in /etc/practable

cp ./files/vna-data /usr/local/bin/
cp ./files/relay-rules /usr/local/bin/
cp ./files/odroid/vna /usr/local/bin
cp ./files/odroid/relay /usr/local/bin 
cp ./services/* /etc/systemd/system

#programme the arduino
curl -fsSL https://raw.githubusercontent.com/arduino/arduino-cli/master/install.sh | sh
arduino-cli core update-index
arduino-cli core install arduino:avr
arduino-cli lib install timerinterrupt
cd ./fw
arduino-cli compile --fqbn arduino:avr:nano RFSwitch/ 
arduino-cli upload --port /dev/ttyUSB0 --fqbn arduino:avr:nano RFSwitch/

# start services
systemctl enable calibration.service
systemctl enable relay/service
systemctl enable relay-rules.service
systemctl enable vna-data.service
systemctl start calibration.service
systemctl start relay.service
systemctl start relay-rules.service
systemctl start vna-data.service







