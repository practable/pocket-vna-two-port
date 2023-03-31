import logging
from concurrent.futures import ThreadPoolExecutor

import grpc
import numpy as np
import os
from skrf.calibration import OnePort, TwelveTerm
from skrf.media import DefinedGammaZ0
import skrf as rf
import warnings

from calibrate_pb2 import CalibrateOnePortResponse,CalibrateTwoPortResponse
from calibrate_pb2_grpc import CalibrateServicer, add_CalibrateServicer_to_server

# For tutorial on grpc with python, see
# https://www.ardanlabs.com/blog/2020/06/python-go-grpc.html

class CalibrateServer(CalibrateServicer):
    def CalibrateOnePort(self, request, context):
        logging.info('calibrate request size: %d', len(request.frequency))
        logging.error('not implemented, returning uncalibrated dut data')
        resp=OutliersResponse(frequency=request.frequency, result=frequest.dut)
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
