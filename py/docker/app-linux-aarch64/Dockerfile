# syntax=docker/dockerfile:1
FROM arm64v8/python:3.7-buster AS base
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y libblas3 liblapack3 liblapack-dev libblas-dev libatlas-base-dev gfortran zlib1g
RUN pip install numpy
RUN pip install scipy
RUN pip install scikit-rf
RUN pip install websocket-client

FROM base AS app
ADD calibration.py /
ADD client.py / 
ENV SESSION_URL ws://172.17.0.1:8888/ws/calibration
CMD [ "python", "./client.py" ]
