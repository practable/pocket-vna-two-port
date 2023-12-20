#!/bin/bash

# we assume jump is already installed
# we assume that relay access and token files are already present in /etc/practable
# assume that /etc/practable/id has the correct id for the experiment, e.g. pvna05

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

#get relay tokens
export FILES=$(/home/odroid/files.link)
export PRACTABLE_ID=$(cat /etc/practable/id)
cd /etc/practable
wget $FILES/st-ed0-data.access.$PRACTABLE_ID -O  st-ed0-data.access
wget $FILES/st-ed0-data.token.$PRACTABLE_ID -O   st-ed0-data.token

# start services
systemctl enable calibration.service
systemctl enable relay/service
systemctl enable relay-rules.service
systemctl enable vna-data.service
systemctl start calibration.service
systemctl start relay.service
systemctl start relay-rules.service
systemctl start vna-data.service







