#!/bin/bash
# usage
# ./connectlocalswitch.sh /dev/ttyUSB1
# asssumes session-relay is running locally, e.g. session-relay host &
websocat ws://localhost:8888/ws/rfswitch tcp-listen:127.0.0.1:9999 --text &
socat $1,echo=0,b57600,crnl tcp:127.0.0.1:9999 &  
# can check with
# websocat ws://localhost:8888/ws/rfswitch -
# {"set":"port","to":"open"} 
# should get {"report":"port","is":"open"}
