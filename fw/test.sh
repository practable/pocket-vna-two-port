#!/bin/sh
# usage: ./test.sh ttyUSB1
websocat ws://localhost:8888/ws/rfswitch tcp-listen:127.0.0.1:9999 --text &
sleep 5s
socat /dev/${1},echo=0,b57600,crnl tcp:127.0.0.1:9999
