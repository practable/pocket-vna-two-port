#!/bin/bash

# reminder to user to modify kernel if 4.9
echo "for kernel 4.9, cgroup error with docker requires kernel mod"
echo "edit /media/boot/boot.ini to include systemd.unified_cgroup_hierarchy=0 (see comments in this script"
# sudo cp /media/boot/boot.ini ./
# sudo nano boot.ini
# sudo diff /media/boot/boot.ini boot.ini
# sudo cp boot.ini /media/boot/boot.ini
# sudo reboot
# setenv bootargs "root=UUID=e139ce78-9841-40fe-8823-96a304a09859 rootwait rw ${condev} ${amlogic} no_console_suspend fsck.repair=yes net.ifnames=0 elevator=noop hdmimode=${hdmimode} cvbsmode=576cvbs max_freq_a55=${max_freq_a55} maxcpus=${maxcpus} voutmode=${voutmode} ${cmode} disablehpd=${disablehpd} cvbscable=${cvbscable} overscan=${overscan} ${hid_quirks} monitor_onoff=${monitor_onoff} logo=osd0,loaded ${cec_enable} sdrmode=${sdrmode} enable_wol=${enable_wol} systemd.unified_cgroup_hierarchy=0"



# we assume jump is already installed
# we assume that relay access and token files are already present in /etc/practable
# assume that /etc/practable/id has the correct id for the experiment, e.g. pvna05
apt-get install curl
curl -sSL https://get.docker.com | sh
sudo usermod -aG docker $USER 
sudo systemctl enable docker

cp ./files/vna-data /usr/local/bin/
cp ./files/relay-rules /usr/local/bin/
cp ./files/odroid/vna /usr/local/bin
cp ./files/odroid/relay /usr/local/bin 
cp ./services/* /etc/systemd/system
cp ../lib/arm64/libPocketVnaApi.so /usr/lib/libPocketVnaApi.so.0
cp ../lib/arm64/libPocketVnaApi.so /usr/lib/libPocketVnaApi.so.1

#programme the arduino
curl -fsSL https://raw.githubusercontent.com/arduino/arduino-cli/master/install.sh | sh
./bin/arduino-cli core update-index
./bin/arduino-cli core install arduino:avr
./bin/arduino-cli lib install timerinterrupt
./bin/arduino-cli compile --fqbn arduino:avr:nano ../fw/RFSwitch/ 
./bin/arduino-cli upload --port /dev/ttyUSB0 --fqbn arduino:avr:nano ../fw/RFSwitch/

#get relay tokens
export FILES=$(cat /home/odroid/files.link)
export PRACTABLE_ID=$(cat /etc/practable/id)
cd /etc/practable
wget $FILES/st-ed0-data.access.$PRACTABLE_ID -O  st-ed0-data.access
wget $FILES/st-ed0-data.token.$PRACTABLE_ID -O   st-ed0-data.token

# start services
systemctl enable calibration.service
systemctl enable relay.service
systemctl enable relay-rules.service
systemctl enable vna-data.service
systemctl start calibration.service
systemctl start relay.service
systemctl start relay-rules.service
systemctl start vna-data.service







