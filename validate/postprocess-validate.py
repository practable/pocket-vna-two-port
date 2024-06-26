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

# practable results
prac_short = rf.Network('./validate/short-validate.s2p', name="prac_short")
prac_open = rf.Network('./validate/open-validate.s2p', name="prac_open")
prac_load = rf.Network('./validate/load-validate.s2p', name="prac_load")
prac_thru = rf.Network('./validate/thru-validate.s2p', name="prac_thru")
prac_dut1 = rf.Network('./validate/dut1-validate.s2p', name="prac_dut1")
prac_dut2 = rf.Network('./validate/dut2-validate.s2p', name="prac_dut2")
prac_dut3 = rf.Network('./validate/dut3-validate.s2p', name="prac_dut3")
prac_dut4 = rf.Network('./validate/dut4-validate.s2p', name="prac_dut4")

# manual results from Alex using pvna
pvna_short_sw1 = rf.Network('../doc/from_alex/PocketVNA/PVNA_SW1_SHORT.s2p', name="pvna_short_sw1")
pvna_open_sw1 = rf.Network('../doc/from_alex/PocketVNA/PVNA_SW1_OPEN.s2p', name="pvna_open_sw1")
pvna_load_sw1 = rf.Network('../doc/from_alex/PocketVNA/PVNA_SW1_LOAD.s2p', name="pvna_load_sw1")
pvna_short_sw2 = rf.Network('../doc/from_alex/PocketVNA/PVNA_SW2_SHORT.s2p', name="pvna_short_sw2")
pvna_open_sw2 = rf.Network('../doc/from_alex/PocketVNA/PVNA_SW2_OPEN.s2p', name="pvna_open_sw2")
pvna_load_sw2 = rf.Network('../doc/from_alex/PocketVNA/PVNA_SW2_LOAD.s2p', name="pvna_load_sw2")
pvna_thru = rf.Network('../doc/from_alex/PocketVNA/PVNA_THRU.s2p', name="pvna_thru")
pvna_dut1 = rf.Network('../doc/from_alex/PocketVNA/PVNA_DUT_1.s2p', name="pvna_dut1")
pvna_dut2 = rf.Network('../doc/from_alex/PocketVNA/PVNA_DUT_2.s2p', name="pvna_dut2")
pvna_dut3 = rf.Network('../doc/from_alex/PocketVNA/PVNA_DUT_3.s2p', name="pvna_dut3")
pvna_dut4 = rf.Network('../doc/from_alex/PocketVNA/PVNA_DUT_4.s2p', name="pvna_dut4")

# manual results from Alex using Copper Mountain VNA
cm_short_sw1 = rf.Network('../doc/from_alex/CopperMountainVNA/SW1_SHORT.s2p', name="cm_short_sw1")
cm_open_sw1 = rf.Network('../doc/from_alex/CopperMountainVNA/SW1_OPEN.s2p', name="cm_open_sw1")
cm_load_sw1 = rf.Network('../doc/from_alex/CopperMountainVNA/SW1_LOAD.s2p', name="cm_load_sw1")
cm_short_sw2 = rf.Network('../doc/from_alex/CopperMountainVNA/SW2_SHORT.s2p', name="cm_short_sw2")
cm_open_sw2 = rf.Network('../doc/from_alex/CopperMountainVNA/SW2_OPEN.s2p', name="cm_open_sw2")
cm_load_sw2 = rf.Network('../doc/from_alex/CopperMountainVNA/SW2_LOAD.s2p', name="cm_load_sw2")
cm_thru = rf.Network('../doc/from_alex/CopperMountainVNA/THRU.s2p', name="cm_thru")
cm_dut1 = rf.Network('../doc/from_alex/CopperMountainVNA/DUT_1.s2p', name="cm_dut1")
cm_dut2 = rf.Network('../doc/from_alex/CopperMountainVNA/DUT_2.s2p', name="cm_dut2")
cm_dut3 = rf.Network('../doc/from_alex/CopperMountainVNA/DUT_3.s2p', name="cm_dut3")
cm_dut4 = rf.Network('../doc/from_alex/CopperMountainVNA/DUT_4.s2p', name="cm_dut4")


# make a list for ease of applying modifications, e.g. to frequency units
items = []

prac = [prac_short,
prac_open,
prac_load,
prac_thru,
prac_dut1,
prac_dut2,
prac_dut3,
prac_dut4]

pvna = [ pvna_short_sw1 ,
pvna_open_sw1 , 
pvna_load_sw1 , 
pvna_short_sw2 ,
pvna_open_sw2 , 
pvna_load_sw2 , 
pvna_thru , 
pvna_dut1 , 
pvna_dut2 , 
pvna_dut3 , 
pvna_dut4 ]

cm = [ cm_short_sw1 ,
cm_open_sw1 , 
cm_load_sw1 , 
cm_short_sw2 ,
cm_open_sw2 , 
cm_load_sw2 , 
cm_thru , 
cm_dut1 , 
cm_dut2 , 
cm_dut3 , 
cm_dut4 ]

items.extend(prac)
items.extend(pvna)
items.extend(cm)
    
for item in items:
    item.frequency.unit = 'ghz'


plots = {
    "short_sw1_db": {"plots": [{"data":prac_short.plot_s_db,"args":[0,0]},
                  {"data":pvna_short_sw1.plot_s_db,"args":[0,0]},
                  {"data":cm_short_sw1.plot_s_db,"args":[0,0]},
                  ],
                 "ylim":[-10,10],
                  },
    "short_sw2_db": {"plots": [{"data":prac_short.plot_s_db,"args":[1,1]},
                  {"data":pvna_short_sw2.plot_s_db,"args":[1,1]},
                  {"data":cm_short_sw2.plot_s_db,"args":[1,1]},
                  ],
                 "ylim":[-10,10],
                  },  
    "open_sw1_db": {"plots": [{"data":prac_open.plot_s_db,"args":[0,0]},
                  {"data":pvna_open_sw1.plot_s_db,"args":[0,0]},
                  {"data":cm_open_sw1.plot_s_db,"args":[0,0]},
                  ],
                 "ylim":[-10,10],
                  }, 
    "open_sw2_db":{"plots":  [{"data":prac_open.plot_s_db,"args":[1,1]},
                  {"data":pvna_open_sw2.plot_s_db,"args":[1,1]},
                  {"data":cm_open_sw2.plot_s_db,"args":[1,1]},
                  ],
                 "ylim":[-10,10],
                  },    
    "load_sw1_db": {"plots": [{"data":prac_load.plot_s_db,"args":[0,0]},
                  {"data":pvna_load_sw1.plot_s_db,"args":[0,0]},
                  {"data":cm_load_sw1.plot_s_db,"args":[0,0]},
                  ],
                 "ylim":[-60,0],
                  }, 
    "load_sw2_db":{"plots":  [{"data":prac_load.plot_s_db,"args":[1,1]},
                  {"data":pvna_load_sw2.plot_s_db,"args":[1,1]},
                  {"data":cm_load_sw2.plot_s_db,"args":[1,1]},
                  ],
                 "ylim":[-60,0],
                  },      
    "short_sw1_deg": {"plots": [{"data":prac_short.plot_s_deg,"args":[0,0]},
                  {"data":pvna_short_sw1.plot_s_deg,"args":[0,0]},
                  {"data":cm_short_sw1.plot_s_deg,"args":[0,0]},
                  ],
                  "ylim":[-180,180],
                  },
    "short_sw2_deg": {"plots": [{"data":prac_short.plot_s_deg,"args":[1,1]},
                  {"data":pvna_short_sw2.plot_s_deg,"args":[1,1]},
                  {"data":cm_short_sw2.plot_s_deg,"args":[1,1]},
                  ],
                      "ylim":[-180,180],
                  },  
    "open_sw1_deg": {"plots": [{"data":prac_open.plot_s_deg,"args":[0,0]},
                  {"data":pvna_open_sw1.plot_s_deg,"args":[0,0]},
                  {"data":cm_open_sw1.plot_s_deg,"args":[0,0]},
                  ],
                     "ylim":[-180,180],
                  }, 
    "open_sw2_deg":{"plots":  [{"data":prac_open.plot_s_deg,"args":[1,1]},
                  {"data":pvna_open_sw2.plot_s_deg,"args":[1,1]},
                  {"data":cm_open_sw2.plot_s_deg,"args":[1,1]},
                  ],
                    "ylim":[-180,180],
                  },    
    "load_sw1_deg": {"plots": [{"data":prac_load.plot_s_deg,"args":[0,0]},
                  {"data":pvna_load_sw1.plot_s_deg,"args":[0,0]},
                  {"data":cm_load_sw1.plot_s_deg,"args":[0,0]},
                  ],
                     "ylim":[-180,180],
                  }, 
    "load_sw2_deg":{"plots":  [{"data":prac_load.plot_s_deg,"args":[1,1]},
                  {"data":pvna_load_sw2.plot_s_deg,"args":[1,1]},
                  {"data":cm_load_sw2.plot_s_deg,"args":[1,1]},
                  ],
                  },    
    "thru_s11_db": {"plots": [{"data":prac_thru.plot_s_db,"args":[0,0]},
                 {"data":pvna_thru.plot_s_db,"args":[0,0]},
                 {"data":cm_thru.plot_s_db,"args":[0,0]},
                 ],
                "ylim":[-70,10],
                 },   
    "thru_s11_deg": {"plots": [{"data":prac_thru.plot_s_deg,"args":[0,0]},
                {"data":pvna_thru.plot_s_deg,"args":[0,0]},
                {"data":cm_thru.plot_s_deg,"args":[0,0]},
                ],
                "ylim":[-180,180],
                },    
    "thru_s12_db": {"plots": [{"data":prac_thru.plot_s_db,"args":[0,1]},
                 {"data":pvna_thru.plot_s_db,"args":[0,1]},
                 {"data":cm_thru.plot_s_db,"args":[0,1]},
                 ],
                "ylim":[-10,10],
                 },   
    "thru_s12_deg": {"plots": [{"data":prac_thru.plot_s_deg,"args":[0,1]},
                {"data":pvna_thru.plot_s_deg,"args":[0,1]},
                {"data":cm_thru.plot_s_deg,"args":[0,1]},
                ],
                "ylim":[-180,180],
                },      
    "thru_s21_db": {"plots": [{"data":prac_thru.plot_s_db,"args":[1,0]},
                 {"data":pvna_thru.plot_s_db,"args":[1,0]},
                 {"data":cm_thru.plot_s_db,"args":[1,0]},
                 ],
                "ylim":[-10,10],
                 },   
    "thru_s21_deg": {"plots": [{"data":prac_thru.plot_s_deg,"args":[1,0]},
                {"data":pvna_thru.plot_s_deg,"args":[1,0]},
                {"data":cm_thru.plot_s_deg,"args":[1,0]},
                ],
                "ylim":[-180,180],
                },   
    "thru_s22_db": {"plots": [{"data":prac_thru.plot_s_db,"args":[1,1]},
                 {"data":pvna_thru.plot_s_db,"args":[1,1]},
                 {"data":cm_thru.plot_s_db,"args":[1,1]},
                 ],
                "ylim":[-70,10],
                 },   
    "thru_s22_deg": {"plots": [{"data":prac_thru.plot_s_deg,"args":[1,1]},
                {"data":pvna_thru.plot_s_deg,"args":[1,1]},
                {"data":cm_thru.plot_s_deg,"args":[1,1]},
                ],
                "ylim":[-180,180],
                },
    "dut1_s11_db": {"plots": [{"data":prac_dut1.plot_s_db,"args":[0,0]},
                 {"data":pvna_dut1.plot_s_db,"args":[0,0]},
                 {"data":cm_dut1.plot_s_db,"args":[0,0]},
                 ],
                "ylim":[-40,10],
                 },   
    "dut1_s11_deg": {"plots": [{"data":prac_dut1.plot_s_deg,"args":[0,0]},
                {"data":pvna_dut1.plot_s_deg,"args":[0,0]},
                {"data":cm_dut1.plot_s_deg,"args":[0,0]},
                ],
                "ylim":[-180,180],
                },    
    "dut1_s12_db": {"plots": [{"data":prac_dut1.plot_s_db,"args":[0,1]},
                 {"data":pvna_dut1.plot_s_db,"args":[0,1]},
                 {"data":cm_dut1.plot_s_db,"args":[0,1]},
                 ],
                "fylim":[-60,10],
                 },   
    "dut1_s12_deg": {"plots": [{"data":prac_dut1.plot_s_deg,"args":[0,1]},
                {"data":pvna_dut1.plot_s_deg,"args":[0,1]},
                {"data":cm_dut1.plot_s_deg,"args":[0,1]},
                ],
                "ylim":[-180,180],
                },      
    "dut1_s21_db": {"plots": [{"data":prac_dut1.plot_s_db,"args":[1,0]},
                 {"data":pvna_dut1.plot_s_db,"args":[1,0]},
                 {"data":cm_dut1.plot_s_db,"args":[1,0]},
                 ],
                "lim":[-60,10],
                 },   
    "dut1_s21_deg": {"plots": [{"data":prac_dut1.plot_s_deg,"args":[1,0]},
                {"data":pvna_dut1.plot_s_deg,"args":[1,0]},
                {"data":cm_dut1.plot_s_deg,"args":[1,0]},
                ],
                "ylim":[-180,180],
                },   
    "dut1_s22_db": {"plots": [{"data":prac_dut1.plot_s_db,"args":[1,1]},
                 {"data":pvna_dut1.plot_s_db,"args":[1,1]},
                 {"data":cm_dut1.plot_s_db,"args":[1,1]},
                 ],
                "ylim":[-40,10],
                 },   
    "dut1_s22_deg": {"plots": [{"data":prac_dut1.plot_s_deg,"args":[1,1]},
                {"data":pvna_dut1.plot_s_deg,"args":[1,1]},
                {"data":cm_dut1.plot_s_deg,"args":[1,1]},
                ],
                "ylim":[-180,180],
                },
    "dut2_s11_db": {"plots": [{"data":prac_dut2.plot_s_db,"args":[0,0]},
                 {"data":pvna_dut2.plot_s_db,"args":[0,0]},
                 {"data":cm_dut2.plot_s_db,"args":[0,0]},
                 ],
                "ylim":[-40,10],
                 },   
    "dut2_s11_deg": {"plots": [{"data":prac_dut2.plot_s_deg,"args":[0,0]},
                {"data":pvna_dut2.plot_s_deg,"args":[0,0]},
                {"data":cm_dut2.plot_s_deg,"args":[0,0]},
                ],
                "ylim":[-180,180],
                },    
    "dut2_s12_db": {"plots": [{"data":prac_dut2.plot_s_db,"args":[0,1]},
                 {"data":pvna_dut2.plot_s_db,"args":[0,1]},
                 {"data":cm_dut2.plot_s_db,"args":[0,1]},
                 ],
                "ylim":[-30,10],
                 },   
    "dut2_s12_deg": {"plots": [{"data":prac_dut2.plot_s_deg,"args":[0,1]},
                {"data":pvna_dut2.plot_s_deg,"args":[0,1]},
                {"data":cm_dut2.plot_s_deg,"args":[0,1]},
                ],
                "ylim":[-180,180],
                },      
    "dut2_s21_db": {"plots": [{"data":prac_dut2.plot_s_db,"args":[1,0]},
                 {"data":pvna_dut2.plot_s_db,"args":[1,0]},
                 {"data":cm_dut2.plot_s_db,"args":[1,0]},
                 ],
                "ylim":[-30,10],
                 },   
    "dut2_s21_deg": {"plots": [{"data":prac_dut2.plot_s_deg,"args":[1,0]},
                {"data":pvna_dut2.plot_s_deg,"args":[1,0]},
                {"data":cm_dut2.plot_s_deg,"args":[1,0]},
                ],
                "ylim":[-180,180],
                },   
    "dut2_s22_db": {"plots": [{"data":prac_dut2.plot_s_db,"args":[1,1]},
                 {"data":pvna_dut2.plot_s_db,"args":[1,1]},
                 {"data":cm_dut2.plot_s_db,"args":[1,1]},
                 ],
                "ylim":[-40,10],
                 },   
    "dut2_s22_deg": {"plots": [{"data":prac_dut2.plot_s_deg,"args":[1,1]},
                {"data":pvna_dut2.plot_s_deg,"args":[1,1]},
                {"data":cm_dut2.plot_s_deg,"args":[1,1]},
                ],
                "ylim":[-180,180],
                },
    "dut3_s11_db": {"plots": [{"data":prac_dut3.plot_s_db,"args":[0,0]},
                 {"data":pvna_dut3.plot_s_db,"args":[0,0]},
                 {"data":cm_dut3.plot_s_db,"args":[0,0]},
                 ],
                "ylim":[-40,10],
                 },   
    "dut3_s11_deg": {"plots": [{"data":prac_dut3.plot_s_deg,"args":[0,0]},
                {"data":pvna_dut3.plot_s_deg,"args":[0,0]},
                {"data":cm_dut3.plot_s_deg,"args":[0,0]},
                ],
                "ylim":[-180,180],
                },    
    "dut3_s12_db": {"plots": [{"data":prac_dut3.plot_s_db,"args":[0,1]},
                 {"data":pvna_dut3.plot_s_db,"args":[0,1]},
                 {"data":cm_dut3.plot_s_db,"args":[0,1]},
                 ],
                "ylim":[-40,10],
                 },   
    "dut3_s12_deg": {"plots": [{"data":prac_dut3.plot_s_deg,"args":[0,1]},
                {"data":pvna_dut3.plot_s_deg,"args":[0,1]},
                {"data":cm_dut3.plot_s_deg,"args":[0,1]},
                ],
                "ylim":[-180,180],
                },      
    "dut3_s21_db": {"plots": [{"data":prac_dut3.plot_s_db,"args":[1,0]},
                 {"data":pvna_dut3.plot_s_db,"args":[1,0]},
                 {"data":cm_dut3.plot_s_db,"args":[1,0]},
                 ],
                "ylim":[-40,10],
                 },   
    "dut3_s21_deg": {"plots": [{"data":prac_dut3.plot_s_deg,"args":[1,0]},
                {"data":pvna_dut3.plot_s_deg,"args":[1,0]},
                {"data":cm_dut3.plot_s_deg,"args":[1,0]},
                ],
                "ylim":[-180,180],
                },   
    "dut3_s22_db": {"plots": [{"data":prac_dut3.plot_s_db,"args":[1,1]},
                 {"data":pvna_dut3.plot_s_db,"args":[1,1]},
                 {"data":cm_dut3.plot_s_db,"args":[1,1]},
                 ],
                "ylim":[-40,10],
                 },   
    "dut3_s22_deg": {"plots": [{"data":prac_dut3.plot_s_deg,"args":[1,1]},
                {"data":pvna_dut3.plot_s_deg,"args":[1,1]},
                {"data":cm_dut3.plot_s_deg,"args":[1,1]},
                ],
                "ylim":[-180,180],
                },
    "dut4_s11_db": {"plots": [{"data":prac_dut4.plot_s_db,"args":[0,0]},
                 {"data":pvna_dut4.plot_s_db,"args":[0,0]},
                 {"data":cm_dut4.plot_s_db,"args":[0,0]},
                 ],
                "ylim":[-40,10],
                 },   
    "dut4_s11_deg": {"plots": [{"data":prac_dut4.plot_s_deg,"args":[0,0]},
                {"data":pvna_dut4.plot_s_deg,"args":[0,0]},
                {"data":cm_dut4.plot_s_deg,"args":[0,0]},
                ],
                "ylim":[-180,180],
                },    
    "dut4_s12_db": {"plots": [{"data":prac_dut4.plot_s_db,"args":[0,1]},
                 {"data":pvna_dut4.plot_s_db,"args":[0,1]},
                 {"data":cm_dut4.plot_s_db,"args":[0,1]},
                 ],
                "fylim":[-40,10],
                 },   
    "dut4_s12_deg": {"plots": [{"data":prac_dut4.plot_s_deg,"args":[0,1]},
                {"data":pvna_dut4.plot_s_deg,"args":[0,1]},
                {"data":cm_dut4.plot_s_deg,"args":[0,1]},
                ],
                "ylim":[-180,180],
                },      
    "dut4_s21_db": {"plots": [{"data":prac_dut4.plot_s_db,"args":[1,0]},
                 {"data":pvna_dut4.plot_s_db,"args":[1,0]},
                 {"data":cm_dut4.plot_s_db,"args":[1,0]},
                 ],
                "ylim":[-40,10],
                 },   
    "dut4_s21_deg": {"plots": [{"data":prac_dut4.plot_s_deg,"args":[1,0]},
                {"data":pvna_dut4.plot_s_deg,"args":[1,0]},
                {"data":cm_dut4.plot_s_deg,"args":[1,0]},
                ],
                "ylim":[-180,180],
                },   
    "dut4_s22_db": {"plots": [{"data":prac_dut4.plot_s_db,"args":[1,1]},
                 {"data":pvna_dut4.plot_s_db,"args":[1,1]},
                 {"data":cm_dut4.plot_s_db,"args":[1,1]},
                 ],
                "ylim":[-40,10],
                 },   
    "dut4_s22_deg": {"plots": [{"data":prac_dut4.plot_s_deg,"args":[1,1]},
                {"data":pvna_dut4.plot_s_deg,"args":[1,1]},
                {"data":cm_dut4.plot_s_deg,"args":[1,1]},
                ],
                "ylim":[-180,180],
                },      
    }

for name in plots.keys():
    plt.figure()
    for item in plots[name]["plots"]:
        item["data"](item["args"][0], item["args"][1])
    if "ylim" in plots[name]:    
        plt.ylim(plots[name]["ylim"])
    title = name.upper().replace("_"," ")
    title = title.replace("DB","(dB)")
    title = title.replace("DEG","(degrees)")
    plt.title(title)    
    plt.savefig("../img/validate-%s.png"%(name), dpi=300)
    plt.close()    

