#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""

Demo2.py  

demonstrate scikit-rf SOLT TwelveTerm two-port cal (with 10 terms)

@author timothy.d.drysdale@gmail.com

Modified 2022-11-10 from demo.py

"""

import json
import skrf as rf
from skrf.calibration import OnePort, TwelveTerm, SOLT
from skrf.media import DefinedGammaZ0
import matplotlib.pyplot as plt
import numpy as np

# measured files supplied from pocket-VNA measurement
twoport_short = rf.Network('test/measured/twoport-short-p1.s2p', name="short")
twoport_open = rf.Network('test/measured/twoport-open-p1.s2p', name="open")
twoport_load = rf.Network('test/measured/twoport-load-p1.s2p', name="load")
twoport_thru = rf.Network('test/measured/twoport-thru.s2p', name="thru")
twoport_dut  = rf.Network('test/measured/twoport-dut.s2p', name="scikit-rf TwelveTerm") #Name for legend later, not what it is now

measured = [\
            twoport_short,
            twoport_open,
            twoport_load,
            twoport_thru
           ]

f = twoport_short.frequency
standard = DefinedGammaZ0(f)

ideals = [\
        standard.short(nports=2),
        standard.open(nports=2),
        standard.load(1e-99, nports=2), #noreflection Gamma -> 0 (can't be zero, div by zero error)
        standard.thru()
        ]


## create a Calibration instance
cal = TwelveTerm(\
        ideals = ideals,
        measured = measured,
        n_thrus=1,
        )

## run, and apply calibration to a DUT

# run calibration algorithm
cal.run()

# # apply it to a dut
dut_cal = cal.apply_cal(twoport_dut)

# # save results for comparison against automated implementation of this approach
dut_cal.write_touchstone('test/expected/twoport.s2p')

# # check results against supplied data

dut_exp = rf.Network('test/supplied/twoport-dut-cal.s2p', name='supplied by course team')

## Try the 3-term calibration on the same data, for S11

# measured files supplied from pocket-VNA measurement
meas2port = [\
        rf.Network('test/measured/short.s2p'),
        rf.Network('test/measured/open.s2p'),
        rf.Network('test/measured/load.s2p'),
        ]
# the data we want is S11

oneport_short_s11 = rf.Network(frequency=f, s=twoport_short.s[:,0,0], name="short")
oneport_open_s11 = rf.Network(frequency=f, s=twoport_open.s[:,0,0], name="open")
oneport_load_s11 = rf.Network(frequency=f, s=twoport_load.s[:,0,0], name="load")

meas_s11 = [\
         oneport_short_s11,
         oneport_open_s11,
         oneport_load_s11,
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

dut_s11 = rf.Network(frequency=f, s=twoport_dut.s[:,0,0], name="scikit-rf OnePort") #name for legend later, not what it is now
# apply it to a dut
dut_s11_cal = cal_s11.apply_cal(dut_s11)

# save results for comparison against automated implementation of this approach
dut_s11_cal.write_touchstone('test/expected/twoport_s11_OnePort.s1p')

# now check S22
oneport_short_s22 = rf.Network(frequency=f, s=twoport_short.s[:,1,1], name="short")
oneport_open_s22 = rf.Network(frequency=f, s=twoport_open.s[:,1,1], name="open")
oneport_load_s22 = rf.Network(frequency=f, s=twoport_load.s[:,1,1], name="load")

meas_s22 = [\
         oneport_short_s22,
         oneport_open_s22,
         oneport_load_s22,
         ]
## create a Calibration instance
cal_s22 = OnePort(\
        ideals = ideals1,
        measured = meas_s22,
        )
# run calibration algorithm
cal_s22.run()

dut_s22 = rf.Network(frequency=f, s=twoport_dut.s[:,1,1], name="scikit-rf OnePort") #name for legend later, not what it is now
# apply it to a dut
dut_s22_cal = cal_s22.apply_cal(dut_s22)

# check results against supplied data

plt.figure()
plt.title("S21")
dut_cal.plot_s_db(1,0)
dut_exp.plot_s_db(1,0)
plt.savefig("img/twoport-demo-s21-db.png",dpi=300)
plt.show()
plt.close()

plt.figure()
plt.title("S12")
dut_cal.plot_s_db(0,1)
dut_exp.plot_s_db(0,1)
plt.savefig("img/twoport-demo-s12-db.png",dpi=300)
plt.show()
plt.close()

plt.figure()
plt.title("S11")
dut_cal.plot_s_db(0,0)
dut_exp.plot_s_db(0,0)
dut_s11_cal.plot_s_db(0,0)
plt.savefig("img/twoport-demo-s11-db.png",dpi=300)
plt.show()
plt.close()

plt.figure()
plt.title("S22")
dut_cal.plot_s_db(1,1)
dut_exp.plot_s_db(1,1)
dut_s22_cal.plot_s_db(0,0)
plt.savefig("img/twoport-demo-s22-db.png",dpi=300)
plt.show()
plt.close()

plt.figure()
plt.title("S12, S21")
dut_cal.plot_s_db(0,1)
dut_cal.plot_s_db(1,0)
plt.savefig("img/twoport-demo-s12s21-db-twelveterm.png",dpi=300)
plt.show()
plt.close()

plt.figure()
plt.title("S12, S21")
dut_exp.plot_s_db(1,0)
dut_exp.plot_s_db(0,1)
plt.savefig("img/twoport-demo-s12s21-db-supplied.png",dpi=300)
plt.show()
plt.close()


## prep for the JSON DEMO ... do it all again, but with JSON.

def s11re(n):
    return n.s_re[:,0,0].tolist()
def s11im(n):
    return n.s_im[:,0,0].tolist()
def s12re(n):
    return n.s_re[:,0,1].tolist()
def s12im(n):
    return n.s_im[:,0,1].tolist()
def s21re(n):
    return n.s_re[:,1,0].tolist()
def s21im(n):
    return n.s_im[:,1,0].tolist()
def s22re(n):
    return n.s_re[:,1,1].tolist()
def s22im(n):
    return n.s_im[:,1,1].tolist()

request = {
        "cmd":"twoport",
        "freq":f.f.tolist(),
        "short":{
            "s11": {
            "real": s11re(twoport_short),
            "imag": s11im(twoport_short)
            },
            "s12": {
            "real": s12re(twoport_short),
            "imag": s12im(twoport_short)
            },            
            "s21": {
            "real": s21re(twoport_short),
            "imag": s21im(twoport_short)
            },           
            "s22": {
            "real": s22re(twoport_short),
            "imag": s22im(twoport_short)
            }   
        },
          "open": {
                "s11": {
                    "real": s11re(twoport_open),
                    "imag": s11im(twoport_open)
                },
                "s12": {
                    "real": s12re(twoport_open),
                    "imag": s12im(twoport_open)
                },            
                "s21": {
                    "real": s21re(twoport_open),
                    "imag": s21im(twoport_open)
                },           
                "s22": {
                    "real": s22re(twoport_open),
                    "imag": s22im(twoport_open)
                }
          },    
              
          "load": {
               "s11": {
                    "real": s11re(twoport_load),
                    "imag": s11im(twoport_load)
                },
                "s12": {
                    "real": s12re(twoport_load),
                    "imag": s12im(twoport_load)
                },            
                "s21": {
                    "real": s21re(twoport_load),
                    "imag": s21im(twoport_load)
                },           
                "s22": {
                    "real": s22re(twoport_load),
                    "imag": s22im(twoport_load)
                } 
             }, 
            "thru":{
                "s11": {
                    "real": s11re(twoport_thru),
                    "imag": s11im(twoport_thru)
                },
                "s12": {
                    "real": s12re(twoport_thru),
                    "imag": s12im(twoport_thru)
                },            
                "s21": {
                    "real": s21re(twoport_thru),
                    "imag": s21im(twoport_thru)
                },           
                "s22": {
                    "real": s22re(twoport_thru),
                    "imag": s22im(twoport_thru)
                }    
            },
            "dut":{
                "s11": {
                    "real": s11re(twoport_dut),
                    "imag": s11im(twoport_dut)
                },
                "s12": {
                    "real": s12re(twoport_dut),
                    "imag": s12im(twoport_dut)
                },            
                "s21": {
                    "real": s21re(twoport_dut),
                    "imag": s21im(twoport_dut)
                },           
                "s22": {
                    "real": s22re(twoport_dut),
                    "imag": s22im(twoport_dut)
                }    
            }  
        }


with open('test/json/twoport.json', 'w') as f:
     json.dump(request, f)





