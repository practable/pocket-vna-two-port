#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
client.py

websocket client for sending commands and plotting results to validate the
combined system of pocketvna, arduino switch, and python calibration service

@author: timothy.d.drysdale@gmail.com
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
        {"id":"0000","t":0,"cmd":"rc","range":{"start":1000000,"end":4000000000},"size":501,"islog":False,"avg":1,"sparam":{"s11":True,"s12":False,"s21":False,"s22":False}},
        {"id":"0001","t":0,"cmd":"crq","what":"dut","avg":1,"sparam":{"s11":True,"s12":False,"s21":False,"s22":False}},  
        {"id":"0002","t":0,"cmd":"crq","what":"short","avg":1,"sparam":{"s11":True,"s12":False,"s21":False,"s22":False}},  
        {"id":"0003","t":0,"cmd":"crq","what":"open","avg":1,"sparam":{"s11":True,"s12":False,"s21":False,"s22":False}},    
        {"id":"0004","t":0,"cmd":"crq","what":"load","avg":1,"sparam":{"s11":True,"s12":False,"s21":False,"s22":False}}, 
        {"id":"0004","t":0,"cmd":"crq","what":"load","avg":1,"sparam":{"s11":True,"s12":False,"s21":False,"s22":False}},  
        ]

names = ["rr","rc","dut","short","open","load","foo"]
network = []

def resultToNetwork(result, name):
    freq = []
    s11 = []
    
    for r in result:
        freq.append(r["freq"])
        s11.append(r["s11"]["real"] + 1j*r["s11"]["imag"])
        
    return rf.Network(frequency=freq, s=s11, name=name)        
        


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
         plt.xlabel("Frequency (GHz)")
         plt.ylabel("Phase (deg)")
         plt.legend()
         plt.savefig("./validate/%s-deg-unwrap-validate.png"%name, dpi=300)
   
         plt.show()
         plt.close()
         n.write_touchstone(filename="validate/%s-validate.s2p"%name,form="db")
     
         
def on_message(ws, message):
    global step, last_step, command, name, network, names
    # we can exploit the heartbeat to cause on_message to be called
    # then we can send and receive our commands
  
    try:
        obj = json.loads(message)
        
        if validExcludingHeartbeat(obj):
            step = step + 1
        
        #printResult(obj)
        
        if step < len(command):
            plotResult(obj, names[step-1])    
        
        if last_step < step and step < len(command):
            cmd = json.dumps(command[step])
            ws.send(cmd)
            print("step %d: %s"%(step, cmd))
            last_step = step
        
    except Exception as e:
        print(e)
        traceback.print_stack()

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

    
    