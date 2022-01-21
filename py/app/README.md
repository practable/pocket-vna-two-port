## Docker considerations

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



