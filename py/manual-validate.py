#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Tue Nov 22 12:06:56 2022

@author: tim
"""
import json
import numpy as np
import skrf as rf
from skrf.plotting import scale_frequency_ticks, add_markers_to_lines
from skrf.calibration import OnePort, TwelveTerm, EightTerm, SOLT
from skrf.media import DefinedGammaZ0
import matplotlib.pyplot as plt

short2 = rf.Network('test/manual/new-nano/nocal-short-101pts.s2p', name="short")
open2 = rf.Network('test/manual/new-nano/nocal-open-101pts.s2p', name="open")
load2 = rf.Network('test/manual/new-nano/nocal-load-101pts.s2p', name="load")
thru2 = rf.Network('test/manual/new-nano/nocal-thru-101pts.s2p', name="thru")
dut1 = rf.Network('test/manual/new-nano/nocal-dut1-101pts.s2p', name="dut1")
dut2 = rf.Network('test/manual/new-nano/nocal-dut2-101pts.s2p', name="dut2")
dut3 = rf.Network('test/manual/new-nano/nocal-dut3-101pts.s2p', name="dut3")
dut4 = rf.Network('test/manual/new-nano/nocal-dut4-101pts.s2p', name="dut4")

plt.figure()
short2.plot_s_db(0,0)
open2.plot_s_db(0,0)
load2.plot_s_db(0,0)
thru2.plot_s_db(0,0)
dut1.plot_s_db(0,0)
dut2.plot_s_db(0,0)
dut3.plot_s_db(0,0)
dut4.plot_s_db(0,0)
plt.savefig("img/manual/new-nano/port1-s11.png",dpi=300)

plt.figure()
short2.plot_s_db(0,1)
open2.plot_s_db(0,1)
load2.plot_s_db(0,1)
thru2.plot_s_db(0,1)
dut1.plot_s_db(0,1)
dut2.plot_s_db(0,1)
dut3.plot_s_db(0,1)
dut4.plot_s_db(0,1)
plt.savefig("img/manual/new-nano/port1-s12.png",dpi=300)

plt.figure()
short2.plot_s_db(1,0)
open2.plot_s_db(1,0)
load2.plot_s_db(1,0)
thru2.plot_s_db(1,0)
dut1.plot_s_db(1,0)
dut2.plot_s_db(1,0)
dut3.plot_s_db(1,0)
dut4.plot_s_db(1,0)
plt.savefig("img/manual/new-nano/port1-s12.png",dpi=300)

plt.figure()
short2.plot_s_db(1,1)
open2.plot_s_db(1,1)
load2.plot_s_db(1,1)
thru2.plot_s_db(1,1)
dut1.plot_s_db(1,1)
dut2.plot_s_db(1,1)
dut3.plot_s_db(1,1)
dut4.plot_s_db(1,1)
plt.savefig("img/manual/new-nano/port1-s22.png",dpi=300)


"""
Digital multimeter pin values
pin 2 - 1 
pin 3 - 1

Pins 4,5,6 8,9,10
short 100   010
open  011  110   
load  110  001
thru  001  101

These might not be right ....

"""


thru2m = rf.Network('test/manual/supplied-nano/thru.s2p', name="thru")
dut1m = rf.Network('test/manual/supplied-nano/dut1.s2p', name="dut1")
dut2m = rf.Network('test/manual/supplied-nano/dut2.s2p', name="dut2")
dut3m = rf.Network('test/manual/supplied-nano/dut3.s2p', name="dut3")
dut4m = rf.Network('test/manual/supplied-nano/dut4.s2p', name="dut4")

plt.figure()
#short2.plot_s_db(0,0)
#open2.plot_s_db(0,0)
#load2.plot_s_db(0,0)
thru2m.plot_s_db(0,0)
dut1m.plot_s_db(0,0)
dut2m.plot_s_db(0,0)
dut3m.plot_s_db(0,0)
dut4m.plot_s_db(0,0)
plt.savefig("img/manual/supplied-nano/port1-s11.png",dpi=300)

plt.figure()
#short2.plot_s_db(0,1)
#open2.plot_s_db(0,1)
#load2.plot_s_db(0,1)
thru2m.plot_s_db(0,1)
dut1m.plot_s_db(0,1)
dut2m.plot_s_db(0,1)
dut3m.plot_s_db(0,1)
dut4m.plot_s_db(0,1)
plt.savefig("img/manual/supplied-nano/port1-s12.png",dpi=300)

plt.figure()
#short2.plot_s_db(1,0)
#open2.plot_s_db(1,0)
#load2.plot_s_db(1,0)
thru2m.plot_s_db(1,0)
dut1m.plot_s_db(1,0)
dut2m.plot_s_db(1,0)
dut3m.plot_s_db(1,0)
dut4m.plot_s_db(1,0)
plt.savefig("img/manual/supplied-nano/port1-s12.png",dpi=300)

plt.figure()
#short2.plot_s_db(1,1)
#open2.plot_s_db(1,1)
#load2.plot_s_db(1,1)
thru2m.plot_s_db(1,1)
dut1m.plot_s_db(1,1)
dut2m.plot_s_db(1,1)
dut3m.plot_s_db(1,1)
dut4m.plot_s_db(1,1)
plt.savefig("img/manual/supplied-nano/port1-s22.png",dpi=300)



