#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""

Demo.py 

demonstrate scikit-rf SOLT one-port cal

@author timothy.d.drysdale@gmail.com

Created 2022-02-20

"""

import json
import numpy as np
import skrf as rf
import time

import warnings

from skrf.calibration import OnePort
from skrf.media import DefinedGammaZ0

#define keys as variables so that mistypes throw an error
cal_cmd = "cmd"
cal_freq = "freq"
cal_short = "short"
cal_open = "open"
cal_load = "load"
cal_dut= "dut"
cal_real = "real"
cal_imag = "imag"

required = (cal_cmd,cal_freq,cal_short,cal_open,cal_load,cal_dut)   
params = (cal_short,cal_open,cal_load,cal_dut)
parts = (cal_real, cal_imag)

def load_json(filename):
    with open(filename, 'r') as f:
        return json.load(f)
    

def is_oneport(obj):
    #https://stackoverflow.com/questions/40659212/futurewarning-elementwise-comparison-failed-returning-scalar-but-in-the-futur
    with warnings.catch_warnings():
        warnings.simplefilter(action='ignore', category=FutureWarning)
        if not all(key in obj for key in required):
            raise KeyError('Missing one or more required keys')
        
    if not obj[cal_cmd].lower() == "oneport":
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
    nan_index = np.isnan(obj[cal_freq])
    
    for param in params:
        for part in parts:
            nan_index = nan_index | np.isnan(obj[param][part])
            
    #delete all (frequency,param) tuples which have a nan in any part of them
    ok_index = (~nan_index).tolist()
    for param in params:
        for part in parts:
            obj[param][part] = np.array(obj[param][part])[ok_index].tolist()   
            
    obj[cal_freq] = np.array(obj[cal_freq])[ok_index].tolist() 
        
    is_oneport(obj) #throws exception if array lengths no longer consistent
    
       
    return {
            "freq":  np.array(obj[cal_freq]),
            "short": np.array(obj[cal_short][cal_real]) + 1j * np.array(obj[cal_short][cal_imag]),
            "open":  np.array(obj[cal_open][cal_real]) + 1j * np.array(obj[cal_open][cal_imag]),
            "load":  np.array(obj[cal_load][cal_real]) + 1j * np.array(obj[cal_load][cal_imag]),
            "dut":   np.array(obj[cal_dut][cal_real]) + 1j * np.array(obj[cal_dut][cal_imag]),
            }
  
def test_object(N):
    #make these lists so we can serialise this for writing to file
      return {
    "cmd":"oneport",
    "freq": np.linspace(1e6,100e6,num=N).tolist(),
    "short":{
        "real":np.random.rand(N).tolist(),
        "imag":np.random.rand(N).tolist(),
            },
     "open":{
        "real":np.random.rand(N).tolist(),
        "imag":np.random.rand(N).tolist(),
            },               
     "load":{
        "real":np.random.rand(N).tolist(),
        "imag":np.random.rand(N).tolist(),
            },                 
     "dut":{
        "real":np.random.rand(N).tolist(),
        "imag":np.random.rand(N).tolist(),
            }  
    } 
     
def make_networks(obj):
    
    #create frequency using data points in object
    f = rf.Frequency()
    f.f = obj["freq"]
    
    #measured cal networks
    meas = [
            rf.Network(frequency=f,s=obj["short"],name="meas_short"),
            rf.Network(frequency=f,s=obj["open"],name="meas_open"),
            rf.Network(frequency=f,s=obj["load"],name="meas_load"),
            ]
    # ideal cal networks
    standard = DefinedGammaZ0(f)
 
    ideal = [
            standard.short(),
            standard.open(),
            standard.load(1e-99),
            ]
       
    dut = rf.Network(frequency=f, s=obj["dut"], name="dut")

    return dut, ideal, meas 
     
def apply_cal(dut, ideal, meas):
     cal = OnePort(ideals = ideal, measured = meas)
     cal.run()
     return cal.apply_cal(dut)
 
def time_apply_cal(dut, ideal, meas):
    
     cal = OnePort(ideals = ideal, measured = meas)
     
     time_start = time.time()
     
     cal.run()
     
     time_cal = time.time()
     
     cal.apply_cal(dut)
     
     time_apply = time.time()
     
     return np.diff([time_start, time_cal, time_apply])
 
def make_cal(ideal, meas):

     cal = OnePort(ideals = ideal, measured = meas)

     cal.run()
     
     return cal
 
def use_cal(cal, dut):
    
    return cal.apply_cal(dut)
    
def network_to_result(network):
    return {
           "freq": network.f.tolist(),
           "S11": {
                       "Real": np.squeeze(network.s_re).tolist(),
                       "Imag": np.squeeze(network.s_im).tolist(),
                   }
    }
   
        

if __name__ == "__main__":
    
    obj = clean_oneport(load_json('test/json/oneport.json'))
    
    # do tests:
    #check clean of good object
    try:
        obj = clean_oneport(test_object(10))
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
        obj[cal_cmd]= "foo"
        clean_oneport(obj)
    except ValueError:
        pass #expected
    else:
        assert False, "Did not raise ValueError for cmd being unequal to oneport"    

    # check that array lengths are checked
    try:
        obj= test_object(10)
        obj["dut"]["real"] = obj["dut"]["real"][2:]
        clean_oneport(obj)
    except ValueError:
        pass #expected
    else:
        assert False, "Did not raise ValueError for uneven array lengths"     
    
    #check that freq length is compared to array lengths
    try:
        obj= test_object(10)
        obj["freq"] = obj["freq"][2:]
        clean_oneport(obj)
    except ValueError:
        pass #expected
    else:
        assert False, "Did not raise ValueError for uneven array lengths"  
        
        
    #add some NaNs and check they are cleaned out
    try:
        
        obj = test_object(10)
        
        f = np.array(obj[cal_freq])
        
        obj[cal_freq][4] = float("nan")
        obj[cal_dut][cal_real][7] = float("nan")
        obj = clean_oneport(obj)
        
        expected = np.array(f[[True, True, True, True, False, True, True, False, True, True]]).astype('float32')
        actual = np.array(obj[cal_freq]).astype('float32')
      
        assert np.array_equal(expected, actual) 
        
    except:
        assert False, "Error reading and cleaning oneport with nan"    

    # check that we can make networks, with the correct assignments of values
    # including the correct ideal properties for the ideal networks
    
    N = 10
    
    obj = test_object(N)
    
    clean_obj = clean_oneport(obj)
    
    dut, ideal, meas = make_networks(clean_obj)
    
    assert len(ideal) == 3
    assert len(meas) == 3
    assert np.array_equal(obj[cal_freq], dut.f)
    assert np.array_equal(obj[cal_freq], ideal[0].f)
    assert np.array_equal(obj[cal_freq], meas[2].f)
    assert np.array_equal(clean_obj[cal_short],np.squeeze(meas[0].s))
    assert np.array_equal(clean_obj[cal_short],np.squeeze(meas[0].s))
    assert np.array_equal(clean_obj[cal_open],np.squeeze(meas[1].s))
    assert np.array_equal(clean_obj[cal_load],np.squeeze(meas[2].s))
    assert np.array_equal(np.zeros(N), np.squeeze(ideal[0].s_db))
    assert np.array_equal(np.zeros(N), np.squeeze(ideal[1].s_db))
    assert np.all(np.less_equal(np.squeeze(ideal[2].s_db), np.ones(N) * -1000))
    assert np.array_equal(np.ones(N) * 180, np.squeeze(ideal[0].s_deg))
    assert np.array_equal(np.zeros(N) * -180, np.squeeze(ideal[1].s_deg))
    
    
    
    # check how long it takes to prepare and apply calibration
    
    time_start = time.time()
    
    obj = clean_oneport(load_json('test/json/oneport.json'))
 
    time_load = time.time()
    
    dut, ideal, meas = make_networks(obj)
    
    time_network = time.time()
        
    data = apply_cal(dut, ideal, meas)
    
    time_apply = time.time()
    
    times = [time_start, time_load, time_network, time_apply]
    #print(np.diff(times)) #[0.00594997 0.00570345 0.19580579]
    
    assert np.all(np.less_equal(np.diff(times), [30e-3, 30e-3, 200e-3]))
    
    time_with_cal = time_apply - time_network
    
    times = time_apply_cal(dut, ideal, meas)
    
    #[0.07986879 0.02432394]
    assert np.all(np.less_equal(np.diff(times), [100e-3, 100e-3]))
    #print(times) #[0.06130695 0.02129531]

    time_start = time.time()
    
    cal = make_cal(ideal, meas)
    
    time_cal = time.time()
     
    result = use_cal(cal, dut)
    
    time_result = time.time()
    
    assert np.all(np.less_equal(np.diff(times), [100e-3, 100e-3]))
    #print(np.diff([time_start, time_cal, time_result])) #[0.06557608 0.02065706]
    time_without_cal = time_result - time_cal
    
    speedup = time_with_cal / time_without_cal
    
    print("%.2f X speedup if cache cal (%d ms vs %d ms)"%(speedup, time_without_cal*1000, time_with_cal*1000))
    #4.37 X speedup if have separate cal (20 ms vs 90 ms)
    
    # check the cal result against the one we calculated and manually
    # compared to the matlab version earlier
    expected = rf.Network('test/expected/expected.s1p')
    
    N = len(expected.f)
    
    max_db_error = np.ones(N)*0.1
    
    actual_db_error = np.abs(np.squeeze(expected.s_db) - np.squeeze(data.s_db))
    
    assert np.all(np.less_equal(actual_db_error, max_db_error))
    
    max_deg_error = np.ones(N)
    
    actual_deg_error = np.abs(np.squeeze(expected.s_deg) - np.squeeze(data.s_deg))
    
    assert np.all(np.less_equal(actual_deg_error, max_deg_error))
    
    # check result_to_json
    result = network_to_result(data)
    
    assert np.array_equal(result["freq"], data.f)
    assert np.array_equal(result["S11"]["Real"], np.squeeze(data.s_re))
    assert np.array_equal(result["S11"]["Imag"], np.squeeze(data.s_im))
    
    # make small json file for testing (and to check serialisation)
    obj = clean_oneport(test_object(10))
    dut, ideal, meas = make_networks(obj)
    data = apply_cal(dut, ideal, meas)
    result = network_to_result(data)
    
    with open('test/json/result.json', 'w') as f:
        json.dump(result, f)
    
    # make a small input file for testing websocket interface
    obj = test_object(10)
    with open('test/json/test.json', 'w') as f:
        json.dump(obj, f)
        
    # test the test file
    test = '{"cmd": "oneport", "freq": [1000000.0, 12000000.0, 23000000.0, 34000000.0, 45000000.0, 56000000.0, 67000000.0, 78000000.0, 89000000.0, 100000000.0], "short": {"real": [0.16039036792339345, 0.5976873413214076, 0.8329869561763907, 0.3632685710600716, 0.291974032483566, 0.40373224059079027, 0.31292545994460963, 0.6152865378289137, 0.10048725234864964, 0.37351939769296627], "imag": [0.189163890844127, 0.5724452385390489, 0.6267571329911495, 0.2639295007809461, 0.6593732039775437, 0.7783445696862374, 0.783076487771806, 0.804688644344, 0.5623621494218922, 0.3862616197557537]}, "open": {"real": [0.9319954420272576, 0.23082952748682928, 0.17132885145217236, 0.6263633833844254, 0.3990915928806996, 0.6202114392564174, 0.7902464722778033, 0.5773940295267551, 0.06208212179318873, 0.7113602777345744], "imag": [0.6833593050494419, 0.9585408787572861, 0.4924999609945958, 0.6369321377956018, 0.5001621474258618, 0.3137951208455947, 0.3747666930821534, 0.022789746635661023, 0.6207930270913493, 0.8046747192984681]}, "load": {"real": [0.5130937226468185, 0.8263480492152019, 0.6249310064069589, 0.7219952766605156, 0.9806134950633044, 0.9732456064297136, 0.33481221535955774, 0.963575339952412, 0.3153019581201433, 0.1455041936747512], "imag": [0.021596161696429417, 0.27190521527394673, 0.5143147186790451, 0.7469481866947372, 0.8968161666645258, 0.5244011364805536, 0.32627405576132407, 0.6113578500659708, 0.9560667800301624, 0.8819825508704627]}, "dut": {"real": [0.3846485029200616, 0.9866224144218543, 0.22520530189759902, 0.42007624963627843, 0.492184292394408, 0.8807271109985295, 0.02996227785649319, 0.9300122489501913, 0.8650807033211915, 0.1842493672879949], "imag": [0.2291021720232519, 0.6050682338747037, 0.8239899826080297, 0.10248054661636785, 0.9935036057620639, 0.6697206261228016, 0.9932594323400147, 0.3647521295447481, 0.6899004876103721, 0.6561124167518747]}}'
    clean_oneport(json.loads(test))