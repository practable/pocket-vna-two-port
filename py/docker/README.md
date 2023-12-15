# Docker deployment

## Building an image for the raspberry pi (arm32v7)

This needs to be built on the raspberry pi, due to a lack of suitable emulators for the raspberry pi 4. Some earlier models are available with QEMU but not suitable for this task.

### Installing Docker on the RPi4

Ensure you are using a raspberry pi with the latest raspbian pi OS lite 32 (currently kernal 5.15)

```
$uname -a
Linux raspberrypi 5.15.84-v7l+ #1613 SMP Thu Jan 5 12:01:26 GMT 2023 armv7l GNU/Linux
```
Update as usual

```
sudo apt update
sudo apt upgrade
```

Install docker....

```
curl -sSL https://get.docker.com | sh
```
avoid needing sudo to run docker commands
```
sudo usermod -aG docker $USER 
```

After installing docker, enable it to start on statup.
`sudo systemctl enable docker`


log out, log in and try hello world
```
docker run hello-world 
```

```
#get sources
cd ~
mkdir sources
cd sources
git clone https://github.com/practable/pocket-vna-two-port.git
cd pocket-vna-two-port/py/docker/app-linux-arm32v7 
sudo build #docker throws a permission denied if you try to run it as an unprivileged user
```

The build is quite slow compared to the desktop build - grab a cuppa / go for a walk.
build start 16:28

## Notes from first version

The Docker container for linux/amd64 and linux/arm32v7 differ because of the lack of pre-compiled python packages on the arm32v7 (Raspberry Pi 4) platform.

A <600MB image is created on ubuntu desktop, in a few minutes, with this Dockerfile:
```
FROM python:3.7-slim
ADD calibration.py /
ADD client.py /
ENV WEBSOCKET ws://172.17.0.1:8888/ws/calibration
COPY ./requirements.txt /var/www/requirements.txt
RUN pip install -r /var/www/requirements.txt
CMD [ "python", "./client.py" ]
```

The build is much slower on rpi4 than desktop because the python dependencies are being compiled. A ~1.6GB image is created in several hours with this Dockerfile:

```
# syntax=docker/dockerfile:1
FROM arm32v7/python:3.7-buster AS base
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y libblas3 liblapack3 liblapack-dev libblas-dev libatlas-base-dev gfortran zlib1g
RUN pip install numpy
RUN pip install scipy
RUN pip install scikit-rf
RUN pip install websocket-client

FROM base AS app
ADD calibration.py /
ADD client.py / 
ENV WEBSOCKET ws://172.17.0.1:8888/ws/calibration
CMD [ "python", "./client.py" ]
```

We use a two-stage build for the raspberry pi to make updates to the python code quicker to deploy, because it avoids the need to recompile any of the python libraries. 

Note: if you just want to build the base image alone, for another project, this Dockerfile will create a layer with the same hash as the `base` layer in the previous dockerfile. 

```
# syntax=docker/dockerfile:1
FROM arm32v7/python:3.7-buster
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y libblas3 liblapack3 liblapack-dev libblas-dev libatlas-base-dev gfortran zlib1g
RUN pip install numpy
RUN pip install scipy
RUN pip install scikit-rf
RUN pip install websocket-client
```



## Docker considerations on linux-amd64

Alpine base images are very slow to build (x50) because they compile from source, so use slim instead (deb based).

`172.17.0.1` is the ip address of the host, where `session host` will be running

The target can be changed via the environment variable passed to the docker container, but defaults to:

`ws://172.17.0.1:8888/ws/calibration`

To test the service, run two websocat instances (each in a different terminal, one to observe, one to write requests)

To observe:
`websocat ws://localhost:8888/ws/calibration -`

To push a json-format request 
`websocat ws://localhost:8888/ws/calibration readfile:./test/json/oneport.json -B 999999`
or a smaller one (easier to see the request and result without scrolling)
`websocat ws://localhost:8888/ws/calibration readfile:./test/json/test.json`


## Raspberry pi build extra notes

See main README for docker info....


Note: you might see this error when building on rpi4/rpios(buster)

```
Fatal Python error: _Py_InitializeMainInterpreter: can't initialize time
PermissionError: [Errno 1] Operation not permitted
```

An issue with building `time` on rpios (`lsb_release -a` reports I am using buster) is discussed [here](https://community.home-assistant.io/t/migration-to-2021-7-fails-fatal-python-error-init-interp-main-cant-initialize-time/320648/9) and solves the problem:
```
sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 04EE7237B7D453EC 648ACFD622F3D138
echo "deb http://deb.debian.org/debian buster-backports main" | sudo tee -a /etc/apt/sources.list.d/buster-backports.list
sudo apt update
sudo apt install -t buster-backports libseccomp2
```

numerical libraries are required for [building scipy](https://docs.scipy.org/doc/scipy/reference/building/linux.html)
```
sudo apt-get install gcc gfortran python3-dev libopenblas-dev liblapack-dev
```
sudo apt install libblas3 liblapack3 liblapack-dev libblas-dev libatlas-base-dev gfortran zlib

An issue with building numpy on rpi (many error messages produced) is solved by adding buildtools to the Dockerfile [as described here](https://stackoverflow.com/questions/63971185/unable-to-install-numpy-on-docker-python3-7-slim-in-a-raspberry-pi)
Change
```
RUN pip install -r /var/www/requirements.txt
```
to

```
RUN apt-get update && \
    apt-get install -y \
        build-essential \
        make \
        gcc \
    && pip install -r /var/www/requirements.txt \
    && apt-get remove -y --purge make gcc build-essential \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/*
```

The build is much slower on rpi4 than desktop because the python dependencies are being compiled.


This Dockerfile now works on rpi4:
```
FROM arm32v7/python:3.7-buster
ADD calibration.py /
ADD client.py /
ENV WEBSOCKET ws://localhost:8888/ws/calibration
RUN pip install numpy
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y libblas3 liblapack3 liblapack-dev libblas-dev libatlas-base-dev gfortran zlib1g
RUN pip install scipy
RUN pip install scikit-rf
RUN pip install websocket-client
CMD [ "python", "./client.py" ]
```
```
docker tag calibration:latest practable/calibration:arm32v7-3.7-buster-0.1
docker push practable/calibration:arm32v7-3.7-buster-0.1
```

This Dockerfile works but requires a few hours to recompile all the python libraries each time you change the contents of the python code in your app, hence the change to the two-stage build process described at the start of the document.


## Systemd

After installing docker, enable it to start on statup.
`sudo systemctl enable docker`

