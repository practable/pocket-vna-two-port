#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Tue Jan 25 19:57:25 2022

@author: tim
"""

import skrf as rf
import matplotlib.pyplot as plt
f_div=1e9
short2port = rf.Network('test/short_after_cal_v3.s2p')
short1port = rf.Network(frequency = short2port.f/f_div, s=short2port.s[:,0,0], name='short')

open2port = rf.Network('test/open_after_cal_v3.s2p')
open1port = rf.Network(frequency = open2port.f/f_div, s=open2port.s[:,0,0], name='open')

load2port = rf.Network('test/load_after_cal_v3.s2p')
load1port = rf.Network(frequency = load2port.f/f_div, s=load2port.s[:,0,0], name='load')

dut2port = rf.Network('test/dut_after_cal_v3.s2p')
dut1port = rf.Network(frequency = dut2port.f/f_div, s=dut2port.s[:,0,0], name='dut')


plt.figure()
short1port.plot_s_db()
open1port.plot_s_db()
load1port.plot_s_db()
dut1port.plot_s_db()
plt.savefig('SOLD_db.png',dpi=300)
plt.show()

plt.figure()
short1port.plot_s_deg()
open1port.plot_s_deg()
load1port.plot_s_deg()
dut1port.plot_s_deg()
plt.savefig('SOLD_deg.png',dpi=300)
plt.show()

plt.figure()
dut1port.plot_it_all()
plt.savefig('dut_all.png',dpi=300)
plt.show()
