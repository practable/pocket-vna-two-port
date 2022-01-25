#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Tue Jan 25 19:57:25 2022

@author: tim
"""

import skrf as rf
import matplotlib.pyplot as plt

dut2port = rf.Network('test/dut_4ghz_501_slow.s2p')
dut1port = rf.Network(frequency = dut2port.f, s=dut2port.s[:,0,0], name='dut')

exp2port = rf.Network('test/supplied/120mm pocket switch.s2p')
exp1port = rf.Network(frequency = exp2port.f, s=exp2port.s[:,0,0], name='expected')

plt.figure()
dut1port.plot_s_db()
exp1port.plot_s_db()
plt.savefig('DUT_db.png',dpi=300)
plt.show()

plt.figure()
dut1port.plot_s_deg()
exp1port.plot_s_deg()
plt.savefig('DUT_deg.png',dpi=300)
plt.show()

