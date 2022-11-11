# syntax=docker/dockerfile:1
FROM arm32v7/python:3.7-buster
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y libblas3 liblapack3 liblapack-dev libblas-dev libatlas-base-dev gfortran zlib1g
RUN pip install numpy
RUN pip install scipy
RUN pip install scikit-rf
RUN pip install websocket-client
RUN pip install matplotlib
