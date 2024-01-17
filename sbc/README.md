# Configuration information

## Ansible scripts

DO NOT USE

These have not yet been updated to the `grpc` version of the firmware - but are kept as a handy reference for how to e.g. program the arduino 

## Manual installation

To aid rollout to the odroids in Dec 2023, `arm64` versions of `jump`, `relay`, `vna` were compiled on `pvna10` which is running kernel 4.9

```
Linux odroid 4.9.337-13 #1 SMP PREEMPT Tue Nov 28 16:28:39 UTC 2023 aarch64 aarch64 aarch64 GNU/Linux
```

The library for the pocketVNA was prepared for cross-compiled for 4.9 but tested on 5.15 ok.


## Credentials

The `jump` and `relay-rules` services will need tokens that are specific to the practable instance hosting them, so generation scripts are not provided here. They will be provided by the adminstrator of the system.

### Jump

The included jump.service requires a `jump.env` file looks like this example:

```
JUMP_HOST_LOCAL_PORT=22
JUMP_HOST_DEVELOPMENT=true
JUMP_HOST_ACCESS=https://app.practable.io/ed0/jump/api/v1/connect/pvna10
JUMP_HOST_TOKEN=ey...
```

The token has been redacted. 

### Relay-rules

There are two files required for the `relay-rules` service that look this these examples:

`st-ed0-data.access`:
```
https://app.practable.io/ed0/access/session/pvna10-st-data
```

`st-ed0-data.token`:
```
ey....
```

