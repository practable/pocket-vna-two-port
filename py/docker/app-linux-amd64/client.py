#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
client.py

websocket client for calculating calibrations

@author: timothy.d.drysdale@gmail.com
"""

from calibration import * 
import json
import os
import _thread
import time
import traceback
import websocket
debug = True

def on_message(ws, message):
    


    try:
        obj = json.loads(message)
        
        if get_cmd(obj)=="oneport":
        
            obj = clean_oneport(obj)
            dut, ideal, meas = make_networks(obj)
            calibrated = apply_cal(dut, ideal, meas)
            result = network_to_result(calibrated)
            ws.send(json.dumps(result)) 
            
        elif get_cmd(obj)=="twoport":
            if debug:
                print("Request:")
                print(obj)
            obj = clean_twoport(obj)
            dut, ideal, meas = make_networks2(obj)
            calibrated = apply_cal2(dut, ideal, meas)
            result = network_to_result2(calibrated)
            ws.send(json.dumps(result)) 
            if debug:
                print("")
                print("Result:")
                print(json.dumps(result))
            
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
    print("To test:")
    print("$ session relay")
    print("$ websocat ws://localhost:8888/ws/calibration readfile:./test/json/oneport.json -B 999999")
    
    url = os.environ.get("SESSION_URL","ws://172.17.0.1:8888/ws/calibration")

    websocket.enableTrace(False)
    ws = websocket.WebSocketApp(url,
                              #on_open=on_open,
                              on_message=on_message,
                              on_error=on_error,
                              on_close=on_close)

    ws.run_forever()
    
    