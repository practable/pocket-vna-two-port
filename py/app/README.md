alpine is really slow with python ... compiles stuff from source.

RUN apk update
RUN apk add bash nano make automake gcc g++ subversion openblas openblas-dev python3-dev

RUN apt-get update

