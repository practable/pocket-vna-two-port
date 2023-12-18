FROM arm32v7/python:3.12-bookworm
ADD server.py /
ADD calibrate_pb2.py /
ADD calibrate_pb2_grpc.py /
ADD __init__.py /
ENV CALIBRATE_PORT 9001
COPY ./requirements.txt /var/www/requirements.txt
RUN apt-get update && \
    apt-get install -y \
        build-essential \
		cmake \
        make \
        gcc \
    && pip install -r /var/www/requirements.txt \
    && apt-get remove -y --purge make gcc build-essential \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/*
EXPOSE 9001
CMD [ "python", "./server.py" ]
