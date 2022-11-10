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
from skrf.calibration import OnePort, TwelveTerm
from skrf.media import DefinedGammaZ0
import skrf as rf
import time
import warnings
import matplotlib.pyplot as plt


#define keys as variables so that mistypes throw an error
# keys common to all cals
cal_cmd = "cmd"
cal_freq = "freq"
cal_real = "real"
cal_imag = "imag"
cal_s11 = "s11"
cal_s12 = "s12"
cal_s21 = "s21"
cal_s22 = "s22"

# keys for 1-port cal params
cal_short = "short"
cal_open = "open"
cal_load = "load"
cal_dut= "dut"

# keys for the 2-port cal params
cal_thru="thru"

# common elements
parts = (cal_real, cal_imag)

# for one port
required = (cal_cmd,cal_freq,cal_short,cal_open,cal_load,cal_dut)   
params = (cal_short,cal_open,cal_load,cal_dut)

# for two port
required2 = (cal_cmd,cal_freq,cal_short,cal_open,cal_load,cal_thru,cal_dut) 
params2 = (cal_short,cal_open,cal_load,cal_thru,cal_dut)
sparams2 = (cal_s11, cal_s12, cal_s21, cal_s22)


def load_json(filename):
    with open(filename, 'r') as f:
        return json.load(f)
    

def get_cmd(obj):
    return obj[cal_cmd].lower()

    
def is_oneport(obj):
    #https://stackoverflow.com/questions/40659212/futurewarning-elementwise-comparison-failed-returning-scalar-but-in-the-futur
    with warnings.catch_warnings():
        warnings.simplefilter(action='ignore', category=FutureWarning)
        if not all(key in obj for key in required):
            raise KeyError('Missing one or more required keys')
        
    if not get_cmd(obj) == "oneport":
        raise ValueError("cmd is not oneport")
        
    #check all lengths are consistent    
    N = len(obj["freq"])
    
    for param in params:
        for part in parts:
            if not len(obj[param][part])==N:
                raise ValueError("Inconsistent sized data arrays")
                
def is_twoport(obj):
    #https://stackoverflow.com/questions/40659212/futurewarning-elementwise-comparison-failed-returning-scalar-but-in-the-futur
    with warnings.catch_warnings():
        warnings.simplefilter(action='ignore', category=FutureWarning)
        if not all(key in obj for key in required2):
            raise KeyError('Missing one or more required keys')
        
    if not get_cmd(obj) == "twoport":
        raise ValueError("cmd is not twoport")
        
    #check all lengths are consistent    
    N = len(obj["freq"])
    
    for param in params2:
        for sparam in sparams2:
            for part in parts:
                if not len(obj[param][sparam][part])==N:
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

def clean_twoport(obj):
    
    is_twoport(obj) #throws exception if not a two-port

    #find nans
    nan_index = np.isnan(obj[cal_freq])
    
    for param in params2:
        for sparam in sparams2:
            for part in parts:
                nan_index = nan_index | np.isnan(obj[param][sparam][part])
            
    #delete all (frequency,param) tuples which have a nan in any part of them
    ok_index = (~nan_index).tolist()
    for param in params2:
        for sparam in sparams2:
            for part in parts:
                obj[param][sparam][part] = np.array(obj[param][sparam][part])[ok_index].tolist()   
            
    obj[cal_freq] = np.array(obj[cal_freq])[ok_index].tolist() 
        
    is_twoport(obj) #throws exception if array lengths no longer consistent
      
    return {
            "freq":  np.array(obj[cal_freq]),
            "short":{
                "s11":  np.array(obj[cal_short][cal_s11][cal_real]) + 1j * np.array(obj[cal_short][cal_s11][cal_imag]),   
                "s12":  np.array(obj[cal_short][cal_s12][cal_real]) + 1j * np.array(obj[cal_short][cal_s12][cal_imag]),   
                "s21":  np.array(obj[cal_short][cal_s21][cal_real]) + 1j * np.array(obj[cal_short][cal_s21][cal_imag]),   
                "s22":  np.array(obj[cal_short][cal_s22][cal_real]) + 1j * np.array(obj[cal_short][cal_s22][cal_imag])                   
                },
            "open":{
                "s11":  np.array(obj[cal_open][cal_s11][cal_real]) + 1j * np.array(obj[cal_open][cal_s11][cal_imag]),   
                "s12":  np.array(obj[cal_open][cal_s12][cal_real]) + 1j * np.array(obj[cal_open][cal_s12][cal_imag]),   
                "s21":  np.array(obj[cal_open][cal_s21][cal_real]) + 1j * np.array(obj[cal_open][cal_s21][cal_imag]),   
                "s22":  np.array(obj[cal_open][cal_s22][cal_real]) + 1j * np.array(obj[cal_open][cal_s22][cal_imag])                   
                },
            "load":{
                "s11":  np.array(obj[cal_load][cal_s11][cal_real]) + 1j * np.array(obj[cal_load][cal_s11][cal_imag]),   
                "s12":  np.array(obj[cal_load][cal_s12][cal_real]) + 1j * np.array(obj[cal_load][cal_s12][cal_imag]),   
                "s21":  np.array(obj[cal_load][cal_s21][cal_real]) + 1j * np.array(obj[cal_load][cal_s21][cal_imag]),   
                "s22":  np.array(obj[cal_load][cal_s22][cal_real]) + 1j * np.array(obj[cal_load][cal_s22][cal_imag])                   
                },
            "thru":{
                "s11":  np.array(obj[cal_thru][cal_s11][cal_real]) + 1j * np.array(obj[cal_thru][cal_s11][cal_imag]),   
                "s12":  np.array(obj[cal_thru][cal_s12][cal_real]) + 1j * np.array(obj[cal_thru][cal_s12][cal_imag]),   
                "s21":  np.array(obj[cal_thru][cal_s21][cal_real]) + 1j * np.array(obj[cal_thru][cal_s21][cal_imag]),   
                "s22":  np.array(obj[cal_thru][cal_s22][cal_real]) + 1j * np.array(obj[cal_thru][cal_s22][cal_imag])                   
                },
            "dut":{
                "s11":  np.array(obj[cal_dut][cal_s11][cal_real]) + 1j * np.array(obj[cal_dut][cal_s11][cal_imag]),   
                "s12":  np.array(obj[cal_dut][cal_s12][cal_real]) + 1j * np.array(obj[cal_dut][cal_s12][cal_imag]),   
                "s21":  np.array(obj[cal_dut][cal_s21][cal_real]) + 1j * np.array(obj[cal_dut][cal_s21][cal_imag]),   
                "s22":  np.array(obj[cal_dut][cal_s22][cal_real]) + 1j * np.array(obj[cal_dut][cal_s22][cal_imag])                   
                }
  
           }



  
def test_object(N):
    #use lists so we can serialise for writing to file
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

def test_object2(N):
    #use lists so we can serialise for writing to file
    return {
        "cmd":"twoport",
        "freq": np.linspace(1e6,100e6,num=N).tolist(),
        "short":{
            "s11":{
                "real":np.random.rand(N).tolist(),
                "imag":np.random.rand(N).tolist(),
                    },
             "s12":{
                "real":np.random.rand(N).tolist(),
                "imag":np.random.rand(N).tolist(),
                    },               
             "s21":{
                "real":np.random.rand(N).tolist(),
                "imag":np.random.rand(N).tolist(),
                    },
            "s22":{
                "real":np.random.rand(N).tolist(),
                "imag":np.random.rand(N).tolist(),
                    }
        },
        "open":{
            "s11":{
                "real":np.random.rand(N).tolist(),
                "imag":np.random.rand(N).tolist(),
                    },
             "s12":{
                "real":np.random.rand(N).tolist(),
                "imag":np.random.rand(N).tolist(),
                    },               
             "s21":{
                "real":np.random.rand(N).tolist(),
                "imag":np.random.rand(N).tolist(),
                    },
            "s22":{
                "real":np.random.rand(N).tolist(),
                "imag":np.random.rand(N).tolist(),
                    }
        },
        "load":{
            "s11":{
                "real":np.random.rand(N).tolist(),
                "imag":np.random.rand(N).tolist(),
                    },
             "s12":{
                "real":np.random.rand(N).tolist(),
                "imag":np.random.rand(N).tolist(),
                    },               
             "s21":{
                "real":np.random.rand(N).tolist(),
                "imag":np.random.rand(N).tolist(),
                    },
            "s22":{
                "real":np.random.rand(N).tolist(),
                "imag":np.random.rand(N).tolist(),
                    }
        },
        "thru":{
            "s11":{
                "real":np.random.rand(N).tolist(),
                "imag":np.random.rand(N).tolist(),
                    },
             "s12":{
                "real":np.random.rand(N).tolist(),
                "imag":np.random.rand(N).tolist(),
                    },               
             "s21":{
                "real":np.random.rand(N).tolist(),
                "imag":np.random.rand(N).tolist(),
                    },
            "s22":{
                "real":np.random.rand(N).tolist(),
                "imag":np.random.rand(N).tolist(),
                    }
        },    
        "dut":{
            "s11":{
                "real":np.random.rand(N).tolist(),
                "imag":np.random.rand(N).tolist(),
                    },
             "s12":{
                "real":np.random.rand(N).tolist(),
                "imag":np.random.rand(N).tolist(),
                    },               
             "s21":{
                "real":np.random.rand(N).tolist(),
                "imag":np.random.rand(N).tolist(),
                    },
            "s22":{
                "real":np.random.rand(N).tolist(),
                "imag":np.random.rand(N).tolist(),
                    }
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

def make_networks2(obj):
    
    #create frequency using data points in object
    f = rf.Frequency()
    f.f = obj["freq"]

    sp_short =np.zeros((len(f), 2, 2), dtype=complex)
    sp_short[:,0,0] = obj[cal_short][cal_s11]
    sp_short[:,0,1] = obj[cal_short][cal_s12]  
    sp_short[:,1,0] = obj[cal_short][cal_s21]    
    sp_short[:,1,1] = obj[cal_short][cal_s22]

    sp_open =np.zeros((len(f), 2, 2), dtype=complex)
    sp_open[:,0,0] = obj[cal_open][cal_s11]
    sp_open[:,0,1] = obj[cal_open][cal_s12]  
    sp_open[:,1,0] = obj[cal_open][cal_s21]    
    sp_open[:,1,1] = obj[cal_open][cal_s22]
    
    sp_load =np.zeros((len(f), 2, 2), dtype=complex)
    sp_load[:,0,0] = obj[cal_load][cal_s11]
    sp_load[:,0,1] = obj[cal_load][cal_s12]  
    sp_load[:,1,0] = obj[cal_load][cal_s21]    
    sp_load[:,1,1] = obj[cal_load][cal_s22]

    sp_thru =np.zeros((len(f), 2, 2), dtype=complex)
    sp_thru[:,0,0] = obj[cal_thru][cal_s11]
    sp_thru[:,0,1] = obj[cal_thru][cal_s12]  
    sp_thru[:,1,0] = obj[cal_thru][cal_s21]    
    sp_thru[:,1,1] = obj[cal_thru][cal_s22]
    
    sp_dut =np.zeros((len(f), 2, 2), dtype=complex)
    sp_dut[:,0,0] = obj[cal_dut][cal_s11]
    sp_dut[:,0,1] = obj[cal_dut][cal_s12]  
    sp_dut[:,1,0] = obj[cal_dut][cal_s21]    
    sp_dut[:,1,1] = obj[cal_dut][cal_s22]    
    
    
    #measured cal networks
    meas = [
            rf.Network(frequency=f,s=sp_short,name="meas_short"),
            rf.Network(frequency=f,s=sp_open,name="meas_open"),
            rf.Network(frequency=f,s=sp_load,name="meas_load"),
            rf.Network(frequency=f,s=sp_thru,name="meas_thru"),
            ]
    # ideal cal networks
    standard = DefinedGammaZ0(f)
 
    ideal = [
            standard.short(nports=2),
            standard.open(nports=2),
            standard.load(1e-99, nports=2),
            standard.thru(),
            ]
       
    dut = rf.Network(frequency=f, s=sp_dut, name="dut")

    return dut, ideal, meas

     
def apply_cal(dut, ideal, meas):
     cal = OnePort(ideals = ideal, measured = meas)
     cal.run()
     return cal.apply_cal(dut)

def apply_cal2(dut, ideal, meas):
     
     cal = TwelveTerm(ideals = ideal, measured = meas, n_thrus=1)
     cal.run()
     dut_cal = cal.apply_cal(dut)
     f = meas[0].frequency
     standard = DefinedGammaZ0(f)

     # now we need to fix S11, S22
     meas_s11 = [\
        rf.Network(frequency=f, s=meas[0].s[:,0,0], name="short"),
        rf.Network(frequency=f, s=meas[1].s[:,0,0], name="open"),
        rf.Network(frequency=f, s=meas[2].s[:,0,0], name="load")   
              ]

     ideals1 = [\
             standard.short(),
             standard.open(),
             standard.load(1e-99), #noreflection Gamma -> 0 (can't be zero, div by zero error)
             ]


     ## create a Calibration instance
     cal_s11 = OnePort(\
             ideals = ideals1,
             measured = meas_s11,
             )
     # run calibration algorithm
     cal_s11.run()

     dut_s11 = rf.Network(frequency=f, s=dut.s[:,0,0], name="dut_s11") 
     
     # apply it to dut_s11
     dut_s11_cal = cal_s11.apply_cal(dut_s11)

     # now check S22
     meas_s22 = [\
              rf.Network(frequency=f, s=meas[0].s[:,1,1], name="short"),
              rf.Network(frequency=f, s=meas[1].s[:,1,1], name="open"),
              rf.Network(frequency=f, s=meas[2].s[:,1,1], name="load")
              ]
     ## create a Calibration instance
     cal_s22 = OnePort(\
             ideals = ideals1,
             measured = meas_s22,
             )
     # run calibration algorithm
     cal_s22.run()

     dut_s22 = rf.Network(frequency=f, s=dut.s[:,1,1], name="dut_s22")
     
     # apply it to dut_s22
     dut_s22_cal = cal_s22.apply_cal(dut_s22)
     
     # sub in our three-term results for s11, s22
     dut_cal.s[:,0,0] = dut_s11_cal.s[:,0,0]
     dut_cal.s[:,1,1] = dut_s22_cal.s[:,0,0]
     
     return dut_cal
     
 
def time_apply_cal(dut, ideal, meas):
    
     cal = OnePort(ideals = ideal, measured = meas)
     
     time_start = time.time()
     
     cal.run()
     
     time_cal = time.time()
     
     cal.apply_cal(dut)
     
     time_apply = time.time()
     
     return np.diff([time_start, time_cal, time_apply])

  
def time_apply_cal2(dut, ideal, meas):
    
     cal = TwelveTerm(ideals = ideal, measured = meas, n_thrus=1)
     
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

def make_cal2(ideal, meas):

     cal = TwelveTerm(ideals = ideal, measured = meas, n_thrus=1)

     cal.run()

     f = meas[0].frequency
     standard = DefinedGammaZ0(f)

     # now we need to fix S11, S22
     meas_s11 = [\
        rf.Network(frequency=f, s=meas[0].s[:,0,0], name="short"),
        rf.Network(frequency=f, s=meas[1].s[:,0,0], name="open"),
        rf.Network(frequency=f, s=meas[2].s[:,0,0], name="load")   
              ]

     ideals1 = [\
             standard.short(),
             standard.open(),
             standard.load(1e-99), #noreflection Gamma -> 0 (can't be zero, div by zero error)
             ]


     ## create a Calibration instance
     cal_s11 = OnePort(\
             ideals = ideals1,
             measured = meas_s11,
             )
     # run calibration algorithm
     cal_s11.run()

     # now check S22
     meas_s22 = [\
              rf.Network(frequency=f, s=meas[0].s[:,1,1], name="short"),
              rf.Network(frequency=f, s=meas[1].s[:,1,1], name="open"),
              rf.Network(frequency=f, s=meas[2].s[:,1,1], name="load")
              ]
     ## create a Calibration instance
     cal_s22 = OnePort(\
             ideals = ideals1,
             measured = meas_s22,
             )
     # run calibration algorithm
     cal_s22.run()
   
     return cal, cal_s11, cal_s22

def use_cal(cal, dut):
    
    return cal.apply_cal(dut)
    
def use_cal2(cal,cal_s11, cal_s22, dut):

    # cal for S12, S21    
    dut_cal = cal.apply_cal(dut)
    
    f = dut.frequency
    
    # cal for S11
    dut_s11 = rf.Network(frequency=f, s=dut.s[:,0,0], name="dut_s11") 
    dut_s11_cal = cal_s11.apply_cal(dut_s11)
    
    # cal for S22
    dut_s22 = rf.Network(frequency=f, s=dut.s[:,1,1], name="dut_s22")
    dut_s22_cal = cal_s22.apply_cal(dut_s22)
   
    # combine cal results 
    dut_cal.s[:,0,0] = dut_s11_cal.s[:,0,0]
    dut_cal.s[:,1,1] = dut_s22_cal.s[:,0,0]
       
    return dut_cal   
    
    
def network_to_result(network):
    return {
           "freq": network.f.tolist(),
           "S11": {
                       "Real": np.squeeze(network.s_re).tolist(),
                       "Imag": np.squeeze(network.s_im).tolist(),
                   }
    }

def network_to_result2(network):
    return {
           "freq": network.f.tolist(),
           "S11": {
                       "Real": np.squeeze(network.s_re[:,0,0]).tolist(),
                       "Imag": np.squeeze(network.s_im[:,0,0]).tolist(),
                   },
           "S12": {
                       "Real": np.squeeze(network.s_re[:,0,1]).tolist(),
                       "Imag": np.squeeze(network.s_im[:,0,1]).tolist(),
                   },
           "S21": {
                       "Real": np.squeeze(network.s_re[:,1,0]).tolist(),
                       "Imag": np.squeeze(network.s_im[:,1,0]).tolist(),
                   },
           "S22": {
                       "Real": np.squeeze(network.s_re[:,1,1]).tolist(),
                       "Imag": np.squeeze(network.s_im[:,1,1]).tolist(),
                   },
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
    
    ## TWO PORT TESTS  ##
    
     # do tests:
    #check clean of good object
    try:
        obj = clean_twoport(test_object2(10))
    except:
        assert False, "Error reading and cleaning good oneport"
       
    #check clean of object missing key throws exception
    for key in required:
        try:
            clean_twoport(test_object2(10).pop(key))
        except KeyError:
            pass #expected
        else:
            assert False, "Did not raise KeyError for missing key"%key
            
    # check that cmd is set as oneport
    try:
        obj= test_object2(10)
        obj[cal_cmd]= "foo"
        clean_twoport(obj)
    except ValueError:
        pass #expected
    else:
        assert False, "Did not raise ValueError for cmd being unequal to oneport"    

    # check that array lengths are checked
    try:
        obj= test_object2(10)
        obj["dut"]["s11"]["real"] = obj["dut"]["s11"]["real"][2:]
        clean_twoport(obj)
    except ValueError:
        pass #expected
    else:
        assert False, "Did not raise ValueError for uneven array lengths"     
    
    #check that freq length is compared to array lengths
    try:
        obj= test_object2(10)
        obj["freq"] = obj["freq"][2:]
        clean_twoport(obj)
    except ValueError:
        pass #expected
    else:
        assert False, "Did not raise ValueError for uneven array lengths"  
        
        
    #add some NaNs and check they are cleaned out
    try:
        
        obj = test_object2(10)
        
        f = np.array(obj[cal_freq])
        
        obj[cal_freq][4] = float("nan")
        obj[cal_dut][cal_s11][cal_real][7] = float("nan")
        obj = clean_twoport(obj)
        
        expected = np.array(f[[True, True, True, True, False, True, True, False, True, True]]).astype('float32')
        actual = np.array(obj[cal_freq]).astype('float32')
      
        assert np.array_equal(expected, actual) 
        
    except:
        assert False, "Error reading and cleaning oneport with nan"    

    # check that we can make networks, with the correct assignments of values
    # including the correct ideal properties for the ideal networks
    
    N = 10
    
    obj = test_object2(N)
    
    clean_obj = clean_twoport(obj)
    
    dut, ideal, meas = make_networks2(clean_obj)
    
    assert len(ideal) == 4
    assert len(meas) == 4
    assert np.array_equal(obj[cal_freq], dut.f)
    assert np.array_equal(obj[cal_freq], ideal[0].f)
    assert np.array_equal(obj[cal_freq], meas[2].f)
    assert np.array_equal(clean_obj[cal_short][cal_s11],np.squeeze((meas[0].s)[:,0,0]))
    assert np.array_equal(clean_obj[cal_open][cal_s11],np.squeeze((meas[1].s)[:,0,0]))
    assert np.array_equal(clean_obj[cal_load][cal_s11],np.squeeze((meas[2].s)[:,0,0]))
    assert np.array_equal(np.ones(N) * 180, np.squeeze(ideal[0].s_deg[:,0,0]))
    assert np.array_equal(np.zeros(N) * -180, np.squeeze(ideal[1].s_deg[:,0,0]))
    
    # check how long it takes to prepare and apply calibration
    
    time_start = time.time()
    
    obj = clean_twoport(load_json('test/json/twoport.json'))
 
    time_load = time.time()
    
    dut, ideal, meas = make_networks2(obj)
    
    time_network = time.time()
        
    data = apply_cal2(dut, ideal, meas)
    
    time_apply = time.time()
    
    times = [time_start, time_load, time_network, time_apply]
    
    time_with_cal = time_apply - time_network
    
    times = time_apply_cal2(dut, ideal, meas)
    
    time_start = time.time()
    
    cal, cal_for_s11, cal_for_s22 = make_cal2(ideal, meas)
    
    time_cal = time.time()
     
    result = use_cal2(cal, cal_for_s11, cal_for_s22, dut) #use 2port version
    
    time_result = time.time()
    
    time_without_cal = time_result - time_cal
    
    speedup = time_with_cal / time_without_cal

    print("%.2f X speedup if cache cal2 (%d ms vs %d ms)"%(speedup, time_without_cal*1000, time_with_cal*1000))
    #9.81 X speedup if cache cal2 (36 ms vs 353 ms)
    
    # check the cal result against the one we calculated and manually
    # compared to the matlab version earlier
    dut_exp = rf.Network('test/expected/twoport.s2p', name="validated python demo")
    dut_cal = result
    dut_cal.Name="calibration service"
    
    plt.figure()
    plt.title("S21")
    dut_cal.plot_s_db(1,0)
    dut_exp.plot_s_db(1,0)
    plt.savefig("img/twoport-cal-s21-db.png",dpi=300)
    plt.show()
    plt.close()

    plt.figure()
    plt.title("S12")
    dut_cal.plot_s_db(0,1)
    dut_exp.plot_s_db(0,1)
    plt.savefig("img/twoport-cal-s12-db.png",dpi=300)
    plt.show()
    plt.close()

    plt.figure()
    plt.title("S11")
    dut_cal.plot_s_db(0,0)
    dut_exp.plot_s_db(0,0)
    plt.savefig("img/twoport-cal-s11-db.png",dpi=300)
    plt.show()
    plt.close()

    plt.figure()
    plt.title("S22")
    dut_cal.plot_s_db(1,1)
    dut_exp.plot_s_db(1,1)
    plt.savefig("img/twoport-cal-s22-db.png",dpi=300)
    plt.show()
    plt.close()
    
    N = len(dut_exp.f)
        
    max_db_error = np.ones(N)*0.1
    
    actual_db_error = np.abs(np.squeeze(dut_exp.s_db[:,0,0]) - np.squeeze(dut_cal.s_db[:,0,0]))
    assert np.all(np.less_equal(actual_db_error, max_db_error))
    
    actual_db_error = np.abs(np.squeeze(dut_exp.s_db[:,0,1]) - np.squeeze(dut_cal.s_db[:,0,1]))
    assert np.all(np.less_equal(actual_db_error, max_db_error))

    actual_db_error = np.abs(np.squeeze(dut_exp.s_db[:,1,0]) - np.squeeze(dut_cal.s_db[:,1,0]))
    assert np.all(np.less_equal(actual_db_error, max_db_error))

    actual_db_error = np.abs(np.squeeze(dut_exp.s_db[:,1,1]) - np.squeeze(dut_cal.s_db[:,1,1]))
    assert np.all(np.less_equal(actual_db_error, max_db_error))    
    
    max_deg_error = np.ones(N)
    
    actual_deg_error = np.abs(np.squeeze(dut_exp.s_deg[:,0,0]) - np.squeeze(dut_cal.s_deg[:,0,0]))
    assert np.all(np.less_equal(actual_deg_error, max_deg_error))

    actual_deg_error = np.abs(np.squeeze(dut_exp.s_deg[:,0,1]) - np.squeeze(dut_cal.s_deg[:,0,1]))
    assert np.all(np.less_equal(actual_deg_error, max_deg_error))
    
    actual_deg_error = np.abs(np.squeeze(dut_exp.s_deg[:,1,0]) - np.squeeze(dut_cal.s_deg[:,1,0]))
    assert np.all(np.less_equal(actual_deg_error, max_deg_error))
    
    actual_deg_error = np.abs(np.squeeze(dut_exp.s_deg[:,1,1]) - np.squeeze(dut_cal.s_deg[:,1,1]))
    assert np.all(np.less_equal(actual_deg_error, max_deg_error))
    
    # check result_to_json
    result = network_to_result2(dut_cal)
    
    assert np.array_equal(result["freq"], dut_cal.f)
    assert np.array_equal(result["S11"]["Real"], np.squeeze(dut_cal.s_re[:,0,0]))
    assert np.array_equal(result["S11"]["Imag"], np.squeeze(dut_cal.s_im[:,0,0]))
    assert np.array_equal(result["S12"]["Real"], np.squeeze(dut_cal.s_re[:,0,1]))
    assert np.array_equal(result["S12"]["Imag"], np.squeeze(dut_cal.s_im[:,0,1]))
    assert np.array_equal(result["S21"]["Real"], np.squeeze(dut_cal.s_re[:,1,0]))
    assert np.array_equal(result["S21"]["Imag"], np.squeeze(dut_cal.s_im[:,1,0]))
    assert np.array_equal(result["S22"]["Real"], np.squeeze(dut_cal.s_re[:,1,1]))
    assert np.array_equal(result["S22"]["Imag"], np.squeeze(dut_cal.s_im[:,1,1]))    
    
    
    # make small json file for testing (and to check serialisation)
    obj = clean_twoport(test_object2(10))
    dut, ideal, meas = make_networks2(obj)
    data = apply_cal2(dut, ideal, meas)
    result = network_to_result2(data)
    
    with open('test/json/result2.json', 'w') as f:
        json.dump(result, f)
    
    # make a small input file for testing websocket interface
    obj = test_object2(10) #don't clean it, as must stay list
    with open('test/json/test2.json', 'w') as f:
        json.dump(obj, f)

    # Opening JSON file
    f = open('test/json/test2.json')
    test = json.load(f)  
    # test the test file
    clean_twoport(test)

