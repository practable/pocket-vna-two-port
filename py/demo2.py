#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""

Demo2.py  

demonstrate scikit-rf SOLT TwelveTerm two-port cal (with 10 terms)

@author timothy.d.drysdale@gmail.com

Modified 2022-11-10 from demo.py

"""


import skrf as rf
from skrf.calibration import TwelveTerm
from skrf.media import DefinedGammaZ0
import matplotlib.pyplot as plt
import numpy as np

# measured files supplied from pocket-VNA measurement


sp1 =  rf.Network('test/measured/twoport-short-p1.s2p')
sp2 =  rf.Network('test/measured/twoport-short-p1.s2p')
op1 =  rf.Network('test/measured/twoport-open-p1.s2p')
op2 =  rf.Network('test/measured/twoport-open-p2.s2p')
lp1 =  rf.Network('test/measured/twoport-load-p1.s2p')
lp2 =  rf.Network('test/measured/twoport-load-p2.s2p')
        
f = sp1.frequency

ss =np.zeros((len(f), 2, 2), dtype=complex)
ss[:,0,0] = sp1.s[:,0,0]
ss[:,1,1] = sp2.s[:,0,0]

os =np.zeros((len(f), 2, 2), dtype=complex)
os[:,0,0] = op1.s[:,0,0]
os[:,1,1] = op2.s[:,0,0]

ls =np.zeros((len(f), 2, 2), dtype=complex)
ls[:,0,0] = lp1.s[:,0,0]
ls[:,1,1] = lp2.s[:,0,0]

twoport_short = rf.Network(frequency=f, s=ss, name="short")
twoport_open = rf.Network(frequency=f, s=ss, name="open")
twoport_load = rf.Network(frequency=f, s=ss, name="load")
twoport_thru = rf.Network('test/measured/twoport-thru.s2p', name="thru")
twoport_dut  = rf.Network('test/measured/twoport-dut.s2p', name="dut")

measured = [\
            twoport_short,
            twoport_open,
            twoport_load,
            twoport_thru
           ]

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
# dut2port = rf.Network('test/supplied/DUTuncal.s2p')
# dut1port = rf.Network(frequency=dut2port.frequency, s=dut2port.s[:,0,0], name="scikit cal")
# dut_caled = cal.apply_cal(dut1port)

# # save results for comparison against automated implementation of this approach
# dut_caled.write_touchstone('test/expected/expected.s1p')

# # check results against supplied data

# expected2port = rf.Network('test/supplied/DUTcal.s2p')
# expected1port = rf.Network(frequency=expected2port.frequency, s=expected2port.s[:,0,0], name="matlab cal")

# plt.figure()
# dut_caled.plot_s_db()
# expected1port.plot_s_db()
# plt.savefig("img/demo-db.png",dpi=300)
# plt.show()
# plt.close()

# plt.figure()
# dut_caled.plot_s_deg()
# expected1port.plot_s_deg()
# plt.savefig("img/demo-deg.png",dpi=300)
# plt.show()
# plt.close()

# plt.figure()
# scdb = np.squeeze(dut_caled.s_db)
# mcdb = np.squeeze(expected1port.s_db)
# plt.plot(dut_caled.f, scdb-mcdb)
# plt.xlabel("Frequency (Hz)")
# plt.ylabel("Error (dB)")
# plt.savefig("img/demo-db-error.png",dpi=300)
# plt.show()
# plt.close()

# plt.figure()
# scdeg = np.squeeze(dut_caled.s_deg)
# mcdeg = np.squeeze(expected1port.s_deg)
# plt.plot(dut_caled.f, scdeg-mcdeg)
# plt.ylim([-180,180])
# plt.xlabel("Frequency (Hz)")
# plt.ylabel("Error (deg)")
# plt.savefig("img/demo-deg-error.png",dpi=300)
# plt.show()
# plt.close()


# ## prep for the JSON DEMO ... do it all again, but with JSON.

# # get our arrays out of the network models
# f = meas1port[0].f.tolist()
# mssr = np.squeeze(meas1port[0].s_re).tolist()
# mssi = np.squeeze(meas1port[0].s_im).tolist()

# msor = np.squeeze(meas1port[1].s_re).tolist()
# msoi = np.squeeze(meas1port[1].s_im).tolist()

# mslr = np.squeeze(meas1port[2].s_re).tolist()
# msli = np.squeeze(meas1port[2].s_im).tolist()

# msdr = np.squeeze(dut1port.s_re).tolist()
# msdi = np.squeeze(dut1port.s_im).tolist()


# request = {
#         "cmd":"oneport",
#         "freq":f,
#         "short":{
#             "real":mssr,
#             "imag":mssi
#                 },
#          "open":{
#             "real":msor,
#             "imag":msoi
#                 },               
#          "load":{
#             "real":mslr,
#             "imag":msli
#                 },                 
#          "dut":{
#             "real":msdr,
#             "imag":msdi
#                 }  
#         }

# import json
# with open('test/json/oneport.json', 'w') as f:
#     json.dump(request, f)





