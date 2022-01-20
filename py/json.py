#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""

Demo.py 

demonstrate scikit-rf SOLT one-port cal

@author timothy.d.drysdale@gmail.com

Created 2022-02-20

"""

import json
import matplotlib.pyplot as plt
import numpy as np
import skrf as rf
import warnings

from skrf.calibration import OnePort
from skrf.media import DefinedGammaZ0

#define keys as variables so that mistypes throw an error
cmd = "cmd"
freq = "freq"
cal_short = "short"
cal_open = "open"
cal_load = "load"
dut= "dut"
real = "real"
imag = "imag"

required = (cmd,freq,cal_short,cal_open,cal_load,dut)   
params = (cal_short,cal_open,cal_load,dut)
parts = (real, imag)

def load_json(filename):
    with open(filename, 'r') as f:
        return json.load(f)
    

def is_oneport(obj):
    #https://stackoverflow.com/questions/40659212/futurewarning-elementwise-comparison-failed-returning-scalar-but-in-the-futur
    with warnings.catch_warnings():
        warnings.simplefilter(action='ignore', category=FutureWarning)
        if not all(key in obj for key in required):
            raise KeyError('Missing one or more required keys')
        
    if not obj[cmd].lower() == "oneport":
        raise ValueError("cmd is not oneport")
        
    #check all lengths are consistent    
    N = len(obj["freq"])
    
    for param in params:
        for part in parts:
            if not len(obj[param][part])==N:
                raise ValueError("Inconsistent sized data arrays")
                
                     
        
def clean_oneport(obj):
    
    is_oneport(obj) #throws exception if not a one-port

    #find nans
    nan_index = np.isnan(obj[freq])
    
    for param in params:
        for part in parts:
            nan_index = nan_index | np.isnan(obj[param][part])
            
    #delete all (frequency,param) tuples which have a nan in any part of them
    ok_index = (~nan_index).tolist()
    for param in params:
        for part in parts:
            obj[param][part] = np.array(obj[param][part])[ok_index].tolist()   
            
            
    is_oneport(obj) #throws exception if array lengths no longer consistent
    
       
    return {
            "freq":  obj[freq],
            "short": np.array(obj[cal_short][real]) + 1j * np.array(obj[cal_short][imag]),
            "open":  np.array(obj[cal_open][real]) + 1j * np.array(obj[cal_open][imag]),
            "load":  np.array(obj[cal_load][real]) + 1j * np.array(obj[cal_load][imag]),
            "dut":   np.array(obj[dut][real]) + 1j * np.array(obj[dut][imag]),
            }
  
def test_object(N):
      return {
    "cmd":"oneport",
    "freq": np.linspace(1e6,100e-6,num=N),
    "short":{
        "real":np.random.rand(N),
        "imag":np.random.rand(N),
            },
     "open":{
        "real":np.random.rand(N),
        "imag":np.random.rand(N),
            },               
     "load":{
        "real":np.random.rand(N),
        "imag":np.random.rand(N),
            },                 
     "dut":{
        "real":np.random.rand(N),
        "imag":np.random.rand(N),
            }  
    } 

if __name__ == "__main__":
    
    # do tests:
    #check clean of good object
    try:
        obj = clean_oneport(load_json('test/json/oneport.json'))
    except:
        assert False, "Error reading and cleaning good oneport"
       
    #check clean of object missing key throws exception
    for key in required:
        try:
            clean_oneport(test_object(10).pop(key))
        except KeyError:
            pass #expected
        else:
            assert False, "Did not raise KeyError for missing key"%key
            
    # check that cmd is set as oneport
    try:
        obj= test_object(10)
        obj[cmd]= "foo"
        clean_oneport(obj)
    except ValueError:
        pass #expected
    else:
        assert False, "Did not raise ValueError for cmd being unequal to oneport"    

    #check clean of object missing key
    try:
        clean_oneport(test_object(10).pop(freq))
    except KeyError:
        pass #expected
    else:
        assert False, "Did not raise KeyError for missing key (freq)"  
    


