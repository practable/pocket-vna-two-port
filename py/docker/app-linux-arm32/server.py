import logging
from concurrent.futures import ThreadPoolExecutor

import grpc
import numpy as np
import os
from skrf.calibration import OnePort, TwelveTerm
from skrf.media import DefinedGammaZ0
import skrf as rf

from calibrate_pb2 import CalibrateOnePortResponse,CalibrateTwoPortResponse, SParams, Complex
from calibrate_pb2_grpc import CalibrateServicer, add_CalibrateServicer_to_server

# For tutorial on grpc with python, see
# https://www.ardanlabs.com/blog/2020/06/python-go-grpc.html

def convert_complex_protoc_to_np(points):

    na = np.array([],dtype=np.complex128)

    for point in points:
        na = np.append(na, point.real + 1j * point.imag)
    
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
    
class CalibrateServer(CalibrateServicer):
    def CalibrateOnePort(self, request, context):
        logging.info('CalibrateOnePort request size: %d', len(request.frequency))
        logging.error('not implemented, returning uncalibrated dut data')
        resp=OutliersResponse(frequency=request.frequency, result=request.dut)
        return resp
    def CalibrateTwoPort(self, request, context):
        logging.info('CalibrateTwoPort request size: %d', len(request.frequency))

        #validate inputs are required length
        rl = len(request.frequency)
        items = [
                request.short,
                request.open,
                request.load,
                request.thru,
                request.dut,
            ]
        ll = []
        
        for item in items:
            ll.append(len(item.s11))
            ll.append(len(item.s12))
            ll.append(len(item.s21))
            ll.append(len(item.s22))
        
        for l in ll:
            if not l == rl:
                context.abort(grpc.StatusCode.INVALID_ARGUMENT,"array lengths do not match frequency")
    
              
      
        #create frequency using data points in object
        f = rf.Frequency()
        f.f = request.frequency
        
        # np (numpy) complex format
        np_short = convert_sparams_protoc_to_np(f, request.short)
        np_open  = convert_sparams_protoc_to_np(f, request.open)
        np_load  = convert_sparams_protoc_to_np(f, request.load)
        np_thru  = convert_sparams_protoc_to_np(f, request.thru)
        np_dut   = convert_sparams_protoc_to_np(f, request.dut)
        
        #measured cal networks
        meas = [
                rf.Network(frequency=f,s=np_short,name="meas_short"),
                rf.Network(frequency=f,s=np_open,name="meas_open"),
                rf.Network(frequency=f,s=np_load,name="meas_load"),
                rf.Network(frequency=f,s=np_thru,name="meas_thru"),
                ]
        # ideal cal networks
        standard = DefinedGammaZ0(f)
     
        ideal = [
                standard.short(nports=2),
                standard.open(nports=2),
                standard.load(1e-99, nports=2),
                standard.thru(),
                ]
           
        dut = rf.Network(frequency=f, s=np_dut, name="dut")

        cal = TwelveTerm(ideals = ideal, measured = meas, n_thrus=1)
        cal.run()
        
        result = convert_rf_to_protoc(cal.apply_cal(dut))

        resp=CalibrateTwoPortResponse(frequency=request.frequency, result=result) 
        
        return resp
        
if __name__ == '__main__':
     logging.basicConfig(
          level=logging.INFO,
          format='%(asctime)s - %(levelname)s - %(message)s',
	)
     server = grpc.server(ThreadPoolExecutor())
     add_CalibrateServicer_to_server(CalibrateServer(), server)
     port = os.getenv("CALIBRATE_PORT")
     if port == None:
         port = 9001
     server.add_insecure_port(f'[::]:{port}')
     server.start()
     logging.info('server ready on port %r', port)
     server.wait_for_termination()
