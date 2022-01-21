#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
client.py

websocket client for calculating calibrations

@author: timothy.d.drysdale@gmail.com
"""

from calibration import clean_oneport, make_networks, apply_cal, network_to_result
import json
import os
import _thread
import time
import traceback
import websocket

def on_message(ws, message):
    
   
    try:

        obj = clean_oneport(json.loads(message))
        dut, ideal, meas = make_networks(obj)
        calibrated = apply_cal(dut, ideal, meas)
        result = network_to_result(calibrated)
        ws.send(json.dumps(result))        
        
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
    
    url = os.environ.get("SESSION_URL","ws://172.17.0.1:8888/ws/calibration")

    websocket.enableTrace(False)
    ws = websocket.WebSocketApp(url,
                              on_message=on_message,
                              on_error=on_error,
                              on_close=on_close)

    ws.run_forever()
    
    
