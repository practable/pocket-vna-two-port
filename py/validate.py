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
import matplotlib.pyplot as plt
import numpy as np
import json
import skrf as rf
import _thread
import time
import traceback
import websocket

step = 0
last_step = -1

command = [
        {"cmd":"rr"},
        {"id":"rcal","t":0,"cmd":"rc","range":{"start":1000000,"end":4000000000},"size":3,"islog":False,"avg":1},
        {"id":"shrt","t":0,"cmd":"rq","range":{"start":1000000,"end":4000000000},"size":3,"what":"short","avg":1,"sparam":{"s11":True,"s12":True, "s21":False,"s22":False}}, 
        {"id":"open","t":0,"cmd":"rq","range":{"start":1000000,"end":4000000000},"size":3,"what":"open","avg":1,"sparam":{"s11":True,"s12":True,"s21":False,"s22":False}}, 
        {"id":"load","t":0,"cmd":"rq","range":{"start":1000000,"end":4000000000},"size":3,"what":"load","avg":1,"sparam":{"s11":True,"s12":True,"s21":False,"s22":False}}, 
        {"id":"thru","t":0,"cmd":"rq","range":{"start":1000000,"end":4000000000},"size":3,"what":"thru","avg":1,"sparam":{"s11":True,"s12":True,"s21":False,"s22":False}}, 
        {"id":"dut1","t":0,"cmd":"rq","range":{"start":1000000,"end":4000000000},"size":3,"what":"dut1","avg":1,"sparam":{"s11":True,"s12":True,"s21":False,"s22":False}},     
        {"id":"dut2","t":0,"cmd":"rq","range":{"start":1000000,"end":4000000000},"size":3,"what":"dut2","avg":1,"sparam":{"s11":True,"s12":True,"s21":False,"s22":False}},  
        {"id":"dut3","t":0,"cmd":"rq","range":{"start":1000000,"end":4000000000},"size":3,"what":"dut3","avg":1,"sparam":{"s11":True,"s12":True,"s21":False,"s22":False}},  
        {"id":"dut4","t":0,"cmd":"rq","range":{"start":1000000,"end":4000000000},"size":3,"what":"dut4","avg":1,"sparam":{"s11":True,"s12":True,"s21":False,"s22":False}},  
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
         
         #plt.figure()
         #n.plot_s_deg()
         #plt.savefig("./validate/%s-deg-validate.png"%name, dpi=300)
         #plt.show()
         #plt.close()
         
         # plt.figure()
         # plt.plot(n.f/1e9,np.squeeze(np.unwrap(n.s_deg[:,0,0], period=360)),label="%s, S11"%name)
         # plt.plot(n.f/1e9,np.squeeze(np.unwrap(n.s_deg[:,0,1], period=360)),label="%s, S12"%name)
         # plt.plot(n.f/1e9,np.squeeze(np.unwrap(n.s_deg[:,1,0], period=360)),label="%s, S21"%name)
         # plt.plot(n.f/1e9,np.squeeze(np.unwrap(n.s_deg[:,1,1], period=360)),label="%s, S22"%name)
         # plt.xlabel("Frequency (GHz)")
         # plt.ylabel("Phase (deg)")
         # plt.legend()
         # plt.savefig("./validate/%s-deg-unwrap-validate.png"%name, dpi=300)
         # plt.show()
         # plt.close()
         
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
            print("got valid")
            print(obj)
            plotResult(obj, "foo")  
            
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
    ws = websocket.WebSocketApp("ws://localhost:8888/ws/data",
                              #on_open=on_open,
                              on_message=on_message,
                              on_error=on_error,
                              on_close=on_close)

    ws.run_forever()

    
    