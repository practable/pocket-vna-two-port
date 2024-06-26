#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
validate.py

websocket client for sending commands and plotting results to validate the
combined system of pocketvna, arduino switch, and python calibration service

To run this test:
$ session host &
$ cd pkg/rfswitch
$  ./connectlocalswitch.sh /dev/ttyUSB1 #runs in background
$ cd $REPO/cmd/vna
$ go build
$ export VNA_DESTINATION=ws://localhost:8888/ws/data
$ export VNA_RFSWITCH=ws://localhost:8888/ws/rfswitch
$ export VNA_CALIBRATION=ws://localhost:8888/ws/calibration
$ export VNA_DEVELOPMENT=true #for extra debug messages
$ ./vna unlock
$ ./vna stream #runs in foreground
$ in another terminal
$ cd $REPO/py
$  ./client.py #runs in foreground

Then run this in the IDE if you like

@author: timothy.d.drysdale@gmail.com

two port version Nov 2022
"""
import base64
import json
import matplotlib.pyplot as plt
import numpy as np
import os
import requests
import skrf as rf
import _thread
import time
import traceback
import websocket
from http.client import HTTPConnection 

step = 0
last_step = -1

debug= True

if debug:
    HTTPConnection.debuglevel = 1

command = [
        {"id":"rr","cmd":"rr"},
        {"id":"rcal","t":0,"cmd":"rc","range":{"start":1000000,"end":4000000000},"size":50,"islog":False,"avg":1},
        # {"id":"uncal_shrt","t":0,"cmd":"rq","range":{"start":1000000,"end":4000000000},"size":30,"what":"short","avg":1,"sparam":{"s11":True,"s12":True, "s21":False,"s22":False}}, 
        # {"id":"uncal_open","t":0,"cmd":"rq","range":{"start":1000000,"end":4000000000},"size":30,"what":"open","avg":1,"sparam":{"s11":True,"s12":True,"s21":False,"s22":False}}, 
        # {"id":"uncal_load","t":0,"cmd":"rq","range":{"start":1000000,"end":4000000000},"size":30,"what":"load","avg":1,"sparam":{"s11":True,"s12":True,"s21":False,"s22":False}}, 
        # {"id":"uncal_thru","t":0,"cmd":"rq","range":{"start":1000000,"end":4000000000},"size":30,"what":"thru","avg":1,"sparam":{"s11":True,"s12":True,"s21":False,"s22":False}}, 
        # {"id":"uncal_dut1","t":0,"cmd":"rq","range":{"start":1000000,"end":4000000000},"size":30,"what":"dut1","avg":1,"sparam":{"s11":True,"s12":True,"s21":False,"s22":False}},     
        # {"id":"uncal_dut2","t":0,"cmd":"rq","range":{"start":1000000,"end":4000000000},"size":30,"what":"dut2","avg":1,"sparam":{"s11":True,"s12":True,"s21":False,"s22":False}},  
        # {"id":"uncal_dut3","t":0,"cmd":"rq","range":{"start":1000000,"end":4000000000},"size":30,"what":"dut3","avg":1,"sparam":{"s11":True,"s12":True,"s21":False,"s22":False}},  
        # {"id":"uncal_dut4","t":0,"cmd":"rq","range":{"start":1000000,"end":4000000000},"size":30,"what":"dut4","avg":1,"sparam":{"s11":True,"s12":True,"s21":False,"s22":False}},  
        {"id":"short","t":0,"cmd":"crq","what":"short","avg":1,"sparam":{"s11":True,"s12":True, "s21":True,"s22":True}}, 
        {"id":"open","t":0,"cmd":"crq","what":"open","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}}, 
        {"id":"load","t":0,"cmd":"crq","what":"load","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}}, 
        {"id":"thru","t":0,"cmd":"crq","what":"thru","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}}, 
        {"id":"dut1","t":0,"cmd":"crq","what":"dut1","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}},     
        {"id":"dut2","t":0,"cmd":"crq","what":"dut2","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}},  
        {"id":"dut3","t":0,"cmd":"crq","what":"dut3","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}},  
        {"id":"dut4","t":0,"cmd":"crq","what":"dut4","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}},          
        {"id":"dut1(repeat)","t":0,"cmd":"crq","what":"dut1","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}}, 
        {"id":"dut2(repeat)","t":0,"cmd":"crq","what":"dut2","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}}, 
        {"id":"dut3(repeat)","t":0,"cmd":"crq","what":"dut3","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}}, 
        {"id":"dut4(repeat)","t":0,"cmd":"crq","what":"dut4","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}},
        {"id":"dut3(rev-repeat)","t":0,"cmd":"crq","what":"dut3","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}},          
        {"id":"dut2(rev-repeat)","t":0,"cmd":"crq","what":"dut2","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}}, 
        {"id":"dut1(rev-repeat)","t":0,"cmd":"crq","what":"dut1","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}}, 
        {"id":"dut4(jump-repeat)","t":0,"cmd":"crq","what":"dut4","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}},                
        {"id":"dut1(jump-repeat)","t":0,"cmd":"crq","what":"dut1","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}},        
        {"id":"dut3(jump-repeat)","t":0,"cmd":"crq","what":"dut3","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}},        
        {"id":"dut1(jump-repeat-2)","t":0,"cmd":"crq","what":"dut1","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}},  
        {"id":"dut2(jump-repeat)","t":0,"cmd":"crq","what":"dut2","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}},
        {"id":"dut1(jump-repeat-3)","t":0,"cmd":"crq","what":"dut1","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}},
        {"id":"short(repeat)","t":0,"cmd":"crq","what":"short","avg":1,"sparam":{"s11":True,"s12":True, "s21":True,"s22":True}}, 
        {"id":"open(repeat)","t":0,"cmd":"crq","what":"open","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}}, 
        {"id":"load(repeat)","t":0,"cmd":"crq","what":"load","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}}, 
        {"id":"thru(repeat)","t":0,"cmd":"crq","what":"thru","avg":1,"sparam":{"s11":True,"s12":True,"s21":True,"s22":True}}, 
        ]

network = []

def resultToNetwork(result, name):
    freq = []
    s11 = []
    s12 = []
    s21 = []
    s22 = []
    
    for r in result:
        freq.append(r["freq"])
        s11.append(r["s11"]["real"] + 1j*r["s11"]["imag"])
        s12.append(r["s12"]["real"] + 1j*r["s12"]["imag"])
        s21.append(r["s21"]["real"] + 1j*r["s21"]["imag"])
        s22.append(r["s22"]["real"] + 1j*r["s22"]["imag"])        
    
    s = np.zeros((len(freq), 2, 2), dtype=complex) 
    s[:,0,0] = s11 
    s[:,0,1] = s12
    s[:,1,0] = s21 
    s[:,1,1] = s22    
    return rf.Network(frequency=freq, s=s, name=name)        
        


def validExcludingHeartbeat(obj):
    if "cmd" in obj:
        if (obj["cmd"] != "hb" and obj["cmd"] != ""):
            return True
    return False   
    
def printResult(obj):
     
     if "result" in obj:
         print(obj["result"])
     else:
         print(obj)
   
def plotResult(obj, name):
     if "result" in obj:
         
         n = resultToNetwork(obj["result"], name)
         n.f = n.f / 1e18
         
         plt.figure()
         n.plot_s_db()
         plt.savefig("./validate/%s-db-validate.png"%name, dpi=300)
         plt.show()
         plt.close()
         
         plt.figure()
         n.plot_s_deg()
         plt.savefig("./validate/%s-deg-validate.png"%name, dpi=300)
         plt.show()
         plt.close()
         
         plt.figure()
         plt.plot(n.f/1e9,np.squeeze(np.unwrap(n.s_deg[:,0,0], period=360)),label="%s, S11"%name)
         plt.plot(n.f/1e9,np.squeeze(np.unwrap(n.s_deg[:,0,1], period=360)),label="%s, S12"%name)
         plt.plot(n.f/1e9,np.squeeze(np.unwrap(n.s_deg[:,1,0], period=360)),label="%s, S21"%name)
         plt.plot(n.f/1e9,np.squeeze(np.unwrap(n.s_deg[:,1,1], period=360)),label="%s, S22"%name)
         plt.xlabel("Frequency (GHz)")
         plt.ylabel("Unwrapped phase (deg)")
         plt.legend()
         plt.savefig("./validate/%s-deg-unwrap-validate.png"%name, dpi=300)
         plt.show()
         plt.close()
         
         n.write_touchstone(filename="validate/%s-validate.s2p"%name,form="db")
     else:
         print("No results for ", name," got", obj) 
     
         
def on_message(ws, message):
    global step, last_step, command, name, network, names
    # we can exploit the heartbeat to cause on_message to be called
    # then we can send and receive our commands
  
    try:
        obj = json.loads(message)
        
        if validExcludingHeartbeat(obj):
            #print(obj)
            print(command[step]["id"])
            plotResult(obj, command[step]["id"])  
            
            step = step + 1
            
        if last_step < step and step <= len(command):
            cmd = json.dumps(command[step])
            ws.send(cmd)
            print("step %d: %s"%(step, cmd))
            last_step = step
        
    except IndexError:
        ws.close()

def on_error(ws, error):
    print(error)

def on_close(ws, close_status_code, close_msg):
    print("### closed ###")

def on_open(ws):
    def run(*args):
        for i in range(3):
            time.sleep(1)
            ws.send("Hello %d" % i)
        time.sleep(1)
        ws.close()
        print("thread terminating...")
    _thread.start_new_thread(run, ())

if __name__ == "__main__":

    websocket.enableTrace(False)
 
    local = True
    
    ws_url = "ws://localhost:8888/ws/data"   
    
    if not local:
        with open('../sbc/autogenerated/data.access.pvna05') as f:
            access = str.strip(f.readline()) #remove newline (else token topic error)
        
        with open('../sbc/autogenerated/data.token.pvna05') as f:
            token = str.strip(f.readline()) #remove newline (else invalid header value error)
        
        #b64token = base64.b64encode(str.encode(token))
        headers = {"Authorization": token}
        print(headers)
        print(access)
        resp = requests.post(access, headers = headers)
        
        # This commands works from the command line
        # SESSION_CLIENT_TOKEN=$(cat ~/sources/pocket-vna-two-port/sbc/autogenerated/data.token.pvna05); curl -XPOST https://relay-access.practable.io/session/pvna05-data -H "Authorization: ${SESSION_CLIENT_TOKEN}" -H "Accept application/json"
        # {"uri":"wss://relay.practable.io/session/pvna05-data?code=ddd77176-2c08-4aa7-a2a1-d23fd441ff14"}
   
        if resp.status_code == 200:
            # ws_url = x.URI
            rj = resp.json()
            ws_url = rj["uri"]
            print(ws_url)
        else:
            print(resp.status_code)
            print(resp.content)
            #os._exit(1)            
  
    
    ws = websocket.WebSocketApp(ws_url,
                              #on_open=on_open,
                              on_message=on_message,
                              on_error=on_error,
                              on_close=on_close)

    ws.run_forever()

    
    
