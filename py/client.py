import matplotlib.pyplot as plt
import numpy as np
import grpc
import os
import skrf as rf
import warnings
from calibrate_pb2 import CalibrateTwoPortRequest, SParams, Complex
from calibrate_pb2_grpc import CalibrateStub

def convert_complex_protoc_to_np(points):

    na = np.array([],dtype=np.complex128)

    for point in points:
        na = np.append(na, (point.real + 1j * point.imag))
    
    return na

def convert_complex_np_to_protoc(points):

    pa = []

    for point in points:
        pa.append(Complex(imag=point.imag, real=point.real))

    return pa    


def convert_sparams_protoc_to_np(f, pobj):
    sp =np.zeros((len(f), 2, 2), dtype=complex)
    sp[:,0,0] = convert_complex_protoc_to_np(pobj.s11)
    sp[:,0,1] = convert_complex_protoc_to_np(pobj.s12)
    sp[:,1,0] = convert_complex_protoc_to_np(pobj.s21)
    sp[:,1,1] = convert_complex_protoc_to_np(pobj.s22)
    return sp
     
def convert_rf_to_protoc(rfobj):
    s11 = convert_complex_np_to_protoc(rfobj.s[:,0,0])
    s12 = convert_complex_np_to_protoc(rfobj.s[:,0,1])
    s21 = convert_complex_np_to_protoc(rfobj.s[:,1,0])
    s22 = convert_complex_np_to_protoc(rfobj.s[:,1,1])
    return SParams(s11=s11,s12=s12,s21=s21,s22=s22)

def protoc_from_s2p(filename):
    n = rf.Network(filename)
    return convert_rf_to_protoc(n)

def protoc_from_two_s2p(filename_sw1, filename_sw2):
    n1 = rf.Network(filename_sw1)
    n2 = rf.Network(filename_sw2)
    s = n1.s
    s[:,1,1] = n2.s[:,1,1]
    n1.s = s
    return convert_rf_to_protoc(n1)


if __name__ == "__main__":

    port = os.getenv("CALIBRATE_PORT")
     
    if port == None:
        port = 9001
         
    channel = grpc.insecure_channel(f'localhost:{port}')
    
    stub = CalibrateStub(channel)


    # there should be minimal change if we use calibrated results
    # i.e. we can compare the input and output files to see if they are similar
    # this makes the test more about correct data conversions, as it could
    # easily be passed by a service that just returned the DUT data as it was
    # unless we specifically check for that - not needed at present.
    # files are at: ../doc/CopperMountainVNA

    prefix = '../doc/from_alex/CopperMountainVNA'
    
    cshort = protoc_from_two_s2p(f'{prefix}/SW1_SHORT.s2p',f'{prefix}/SW2_SHORT.s2p')
    copen = protoc_from_two_s2p(f'{prefix}/SW1_OPEN.s2p',f'{prefix}/SW2_OPEN.s2p')
    cload = protoc_from_two_s2p(f'{prefix}/SW1_LOAD.s2p',f'{prefix}/SW2_LOAD.s2p')
    cthru = protoc_from_s2p(f'{prefix}/THRU.s2p')
    cdut = protoc_from_s2p(f'{prefix}/DUT_1.s2p')
    n = rf.Network(f'{prefix}/DUT_1.s2p') #for frequency
    f = n.f
    
    req = CalibrateTwoPortRequest(
        frequency=f,
        short=cshort,
        open=copen,
        load=cload,
        thru=cthru,
        dut=cdut
    )

    resp = stub.CalibrateTwoPort(req)
    s = convert_sparams_protoc_to_np(f,resp.result)
    d = rf.Network(f=f,s=s)
    dr = rf.Network(f'{prefix}/DUT_1.s2p')
    plt.figure()
    markevery=5
    f = f/1e9
    plt.plot(f,d.s_db[:,0,0], "ro", label="s11_cal", markevery=markevery)
    plt.plot(f,d.s_db[:,0,1], "yo", label="s12_cal", markevery=markevery)
    plt.plot(f,d.s_db[:,1,0], "go", label="s21_cal", markevery=markevery)
    plt.plot(f,d.s_db[:,1,1], "bo", label="s22_cal", markevery=markevery)    

    plt.plot(f,dr.s_db[:,0,0], "r:", label="s11_ref", markevery=markevery)
    plt.plot(f,dr.s_db[:,0,1], "y:", label="s12_ref", markevery=markevery)
    plt.plot(f,dr.s_db[:,1,0], "g:", label="s21_ref", markevery=markevery)
    plt.plot(f,dr.s_db[:,1,1], "b:", label="s22_ref", markevery=markevery)  
    plt.xlabel("frequency/GHz")
    plt.ylabel("dB")
    plt.legend()
    plt.savefig("../validate/proto-validate-dut1.png",dpi=300)
    
    
