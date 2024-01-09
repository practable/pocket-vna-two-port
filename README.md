# pocket-vna-two-port

A golang-based web service for the pocketvna.com](https://pocketvna.com/) vector network analyser, with an automated two-port calibration facility provided by additional RF hardware, for use with the open-source [practable.io](https://practable.io] remote lab ecosystem. The user interface is similar to [this](https://github.com/dpreid/pidui), although not yet publically released by the developer (coming soon!)

Status: Prototype has passed measurement quality check, and firmware has been updated to two-ports.

## API

Note that the optional "id" parameter does not affect the command itself, and is provided to help associate replies with the commands that produced them (in case responses come out of order)

### Reasonable range 

```
{"id":"rr","cmd":"rr"}
```
Response:
```
{"id":"rr","t":0,"cmd":"rr","range":{"start":500000,"end":4000000000}}
```

### Calibration

To Perform a full 2-port calibration issue a command like this:

```
{"id":"rcal","t":0,"cmd":"rc","range":{"start":1000000,"end":4000000000},"size":501,"islog":false,"avg":1}
```

501 points is quite slow, so for quick checking functionality, try fewer points.
```
{"id":"rcal","t":0,"cmd":"rc","range":{"start":1000000,"end":4000000000},"size":5,"islog":false,"avg":1}
```

 The `range`, `size`, `isLog` and `avg` are the same as before. This example uses the largest scan size possible (501 points)

### Measurement

These are all the measurements that can be taken (as before, they use the size, and range parameters from the cal):

```
{"id":"short","t":0,"cmd":"crq","what":"short","avg":1,"sparam":{"s11":true,"s12":true, "s21":true,"s22":true}}
{"id":"open","t":0,"cmd":"crq","what":"open","avg":1,"sparam":{"s11":true,"s12":true,"s21":true,"s22":true}}
{"id":"load","t":0,"cmd":"crq","what":"load","avg":1,"sparam":{"s11":true,"s12":true,"s21":true,"s22":true}}
{"id":"thru","t":0,"cmd":"crq","what":"thru","avg":1,"sparam":{"s11":true,"s12":true,"s21":true,"s22":true}}
{"id":"dut1","t":0,"cmd":"crq","what":"dut1","avg":1,"sparam":{"s11":true,"s12":true,"s21":true,"s22":true}}
{"id":"dut2","t":0,"cmd":"crq","what":"dut2","avg":1,"sparam":{"s11":true,"s12":true,"s21":true,"s22":true}}
{"id":"dut3","t":0,"cmd":"crq","what":"dut3","avg":1,"sparam":{"s11":true,"s12":true,"s21":true,"s22":true}}
{"id":"dut4","t":0,"cmd":"crq","what":"dut4","avg":1,"sparam":{"s11":true,"s12":true,"s21":true,"s22":true}} 
```

### Step wise calibration

New feature added to support students learning about the cal process

SC Command (Setup cal, sets the frequency range, number of points, distribution type, and averaging count)
`{"cmd":"sc","range":{"start":1000000,"end":4000000000},"size":20,"islog":false,"avg":1}`

SC Response
```
{"message":"ok","Command":{"id":"","t":0,"cmd":"sc","range":{"start":1000000,"end":4000000000},"size":20,"islog":false,"avg":1,"sparam":{"s11":false,"s12":false,"s21":false,"s22":false},"result":null,"what":"thru"}}
```

MC Short Command:
`{"cmd":"mc","what":"short"}`

MC Short Response (has no data):
```
{"message":"ok","Command":{"id":"","t":0,"cmd":"mc","range":{"start":0,"end":0},"size":0,"islog":false,"avg":0,"sparam":{"s11":false,"s12":false,"s21":false,"s22":false},"result":null,"what":"short"}}
```

MC Open Command:
`{"cmd":"mc","what":"open"}`

MC Open Response (has no data):
```
{"message":"ok","Command":{"id":"","t":0,"cmd":"mc","range":{"start":0,"end":0},"size":0,"islog":false,"avg":0,"sparam":{"s11":false,"s12":false,"s21":false,"s22":false},"result":null,"what":"open"}}
```

MC Load Command:
`{"cmd":"mc","what":"load"}`

MC Load Response (has no data):
```
{"message":"ok","Command":{"id":"","t":0,"cmd":"mc","range":{"start":0,"end":0},"size":0,"islog":false,"avg":0,"sparam":{"s11":false,"s12":false,"s21":false,"s22":false},"result":null,"what":"load"}}
```

MC Thru Command:

`{"cmd":"mc","what":"thru"}`

MC Thru Response (has no data):

```
{"message":"ok","Command":{"id":"","t":0,"cmd":"mc","range":{"start":0,"end":0},"size":0,"islog":false,"avg":0,"sparam":{"s11":false,"s12":false,"s21":false,"s22":false},"result":null,"what":"thru"}}
```

CC Command (cal confirm):

`{"cmd":"cc"}`

CC Command (has data, "what" indicates this is data for the thru with cal applied)
```
{"id":"","t":0,"cmd":"cc","range":{"start":0,"end":0},"size":0,"islog":false,"avg":0,"sparam":{"s11":false,"s12":false,"s21":false,"s22":false},"result":[{"s11":{"real":6.798285684431499e-16,"imag":2.4445757591977986e-15},"s12":{"real":1.0000000000000038,"imag":3.961372257700955e-15},"s21":{"real":0.9999999999999997,"imag":-2.5403542658420376e-15},"s22":{"real":-1.3100837207899224e-16,"imag":1.006857580582699e-15},"freq":1000000},{"s11":{"real":-5.4796192872464335e-18,"imag":3.681361022427108e-16},"s12":{"real":0.9999999999999998,"imag":1.291216167493111e-16},"s21":{"real":0.9999999999999998,"imag":1.0329729339944889e-16},"s22":{"real":-5.483483876630498e-16,"imag":2.274972688902102e-16},"freq":211473684},{"s11":{"real":4.935662660074536e-16,"imag":-7.96132869667085e-17},"s12":{"real":1.0000000000000002,"imag":9.417544964280893e-18},"s21":{"real":1,"imag":3.3389477600632256e-17},"s22":{"real":-5.493611090297354e-17,"imag":2.7110250401774997e-17},"freq":421947368},{"s11":{"real":-3.7997788415861466e-16,"imag":-5.874815142871674e-16},"s12":{"real":0.9999999999999998,"imag":-1.3589387524604561e-17},"s21":{"real":1,"imag":1.358938752460456e-16},"s22":{"real":2.7701351888980784e-16,"imag":2.4761558453947516e-16},"freq":632421052},{"s11":{"real":-1.770339245638541e-16,"imag":1.9778244425109326e-16},"s12":{"real":1,"imag":-8.300789127383292e-17},"s21":{"real":1,"imag":-1.3834648545638821e-17},"s22":{"real":1.8971165707760076e-16,"imag":2.4057070110192795e-16},"freq":842894736},{"s11":{"real":-4.2721621505142654e-16,"imag":-3.6944203527273236e-16},"s12":{"real":1.0000000000000002,"imag":1.0316557900882054e-16},"s21":{"real":1,"imag":2.3212255276984624e-16},"s22":{"real":1.4673180020770904e-16,"imag":-7.563249442097559e-17},"freq":1053368421},{"s11":{"real":-7.434581424175234e-17,"imag":-1.9738283167151276e-16},"s12":{"real":1.0000000000000007,"imag":2.653377447265651e-16},"s21":{"real":1,"imag":1.3668914122277596e-16},"s22":{"real":-2.3401087672749595e-16,"imag":1.0842732839585657e-15},"freq":1263842105},{"s11":{"real":3.6225925572304583e-16,"imag":-6.553331226630913e-16},"s12":{"real":1.0000000000000002,"imag":-7.697793525240971e-18},"s21":{"real":1,"imag":-2.232360122319882e-16},"s22":{"real":-1.664758573382055e-16,"imag":4.426697324102885e-16},"freq":1474315789},{"s11":{"real":-2.1111596138986894e-16,"imag":-3.374000478864848e-16},"s12":{"real":0.9999999999999998,"imag":2.500148396700068e-16},"s21":{"real":0.9999999999999997,"imag":1.388971331500038e-16},"s22":{"real":-9.217893404457907e-16,"imag":5.883978158950974e-16},"freq":1684789473},{"s11":{"real":-9.143566579079719e-17,"imag":-1.211669953147277e-16},"s12":{"real":0.9999999999999994,"imag":-3.384880863400854e-16},"s21":{"real":0.9999999999999997,"imag":-1.9131935314874394e-16},"s22":{"real":7.370339974993562e-16,"imag":-8.28234587808332e-16},"freq":1895263157},{"s11":{"real":6.457509151560786e-16,"imag":3.8407024419460014e-16},"s12":{"real":1,"imag":1.9729520543081967e-16},"s21":{"real":1,"imag":1.1097855305483607e-16},"s22":{"real":6.60122631708177e-16,"imag":1.3745123740607558e-16},"freq":2105736842},{"s11":{"real":-3.3440705904509994e-17,"imag":4.449237964595957e-17},"s12":{"real":1.0000000000000002,"imag":-2.5977695183489855e-17},"s21":{"real":1,"imag":-2.5977695183489855e-17},"s22":{"real":-5.030390119684563e-17,"imag":2.0588438761762348e-16},"freq":2316210526},{"s11":{"real":9.761453389683951e-16,"imag":9.948440742403095e-16},"s12":{"real":1,"imag":3.198724411505024e-18},"s21":{"real":1,"imag":-6.173538114204697e-16},"s22":{"real":-5.60392729997236e-16,"imag":1.011513121527802e-15},"freq":2526684210},{"s11":{"real":3.4281792534273156e-16,"imag":1.8971075667988323e-16},"s12":{"real":0.9999999999999998,"imag":1.1750683925917896e-16},"s21":{"real":1.0000000000000002,"imag":-4.0090568688425766e-16},"s22":{"real":-4.803256717402525e-16,"imag":-1.6853048357543685e-16},"freq":2737157894},{"s11":{"real":7.0222190594856e-16,"imag":8.55812317471566e-16},"s12":{"real":0.9999999999999988,"imag":-1.1169563182435168e-15},"s21":{"real":0.9999999999999999,"imag":9.307969318695973e-17},"s22":{"real":1.2657248559624003e-15,"imag":-4.1238155967451286e-16},"freq":2947631578},{"s11":{"real":3.858821128617499e-16,"imag":-5.908283352150228e-16},"s12":{"real":1.0000000000000004,"imag":2.7559016913199073e-16},"s21":{"real":0.9999999999999999,"imag":-3.1910440636335773e-16},"s22":{"real":-8.995649578270255e-16,"imag":1.3676784350870345e-15},"freq":3158105263},{"s11":{"real":2.1930623628836212e-15,"imag":-4.471183287954572e-16},"s12":{"real":1.0000000000000007,"imag":6.624851261896299e-17},"s21":{"real":1.0000000000000009,"imag":-1.5568400465456302e-15},"s22":{"real":-3.910813152158192e-16,"imag":1.497689622024306e-16},"freq":3368578947},{"s11":{"real":-2.146668654496036e-15,"imag":-1.183624872561583e-15},"s12":{"real":1.0000000000000009,"imag":8.946818283441143e-16},"s21":{"real":0.9999999999999997,"imag":7.437272811749428e-16},"s22":{"real":-7.642394840158961e-16,"imag":1.3994328139316487e-15},"freq":3579052631},{"s11":{"real":2.7529090287958957e-17,"imag":7.844753059553777e-16},"s12":{"real":1.0000000000000007,"imag":4.241877296255027e-16},"s21":{"real":1.0000000000000007,"imag":-1.4812904844065175e-16},"s22":{"real":1.4749059628833424e-17,"imag":7.703227823525666e-16},"freq":3789526315},{"s11":{"real":2.1613952847620654e-15,"imag":-1.4106120953151985e-15},"s12":{"real":1.0000000000000018,"imag":1.8138127619939417e-16},"s21":{"real":0.9999999999999996,"imag":-5.654828022686995e-16},"s22":{"real":8.233651548249787e-16,"imag":1.6094852797772485e-15},"freq":4000000000}],"what":"thru"}
```

 If you try to confirm the cal or take a measurement before completing the cal, then you will get an error message, e.g. sending cc before doing anything else ...
 
```
{"message":"not calibrated yet","Command":{"id":"","t":0,"cmd":"cc","range":{"start":0,"end":0},"size":0,"islog":false,"avg":0,"sparam":{"s11":false,"s12":false,"s21":false,"s22":false},"result":null,"what":""}}
```

 If you have started but not finished the cal, then there will be a hint about what is missing. It's not very sophisticated, it checks for short, open, load thru in that order but it is ok to do the measurements in any order. E.g. if you have only done the `sc` and then `mc` for the `short`, and then try `cc`, you get
```
{"message":"calibration not complete (missing open, maybe others)","Command":{"id":"","t":0,"cmd":"cc","range":{"start":0,"end":0},"size":0,"islog":false,"avg":0,"sparam":{"s11":false,"s12":false,"s21":false,"s22":false},"result":null,"what":""}}
```
But if you try a crq, you don't; get any hints, just a message it is not calibrated yet:

CRQ command

```
{"id":"short","t":0,"cmd":"crq","what":"short","avg":1,"sparam":{"s11":true,"s12":true, "s21":true,"s22":true}}
```

CRQ response

```
{"message":"not calibrated yet","Command":{"id":"short","t":0,"cmd":"crq","what":"short","avg":1,"sparam":{"s11":true,"s12":true,"s21":true,"s22":true},"result":null}}
```

### Trying it out

use the dev-debug because it has been modified to display the data stream more conveniently

request rr:
```
{"id":"rr","t":0,"cmd":"rr","range":{"start":500000,"end":4000000000}}
```
response
```
{"id":"rcal","t":0,"cmd":"rc","range":{"start":1000000,"end":4000000000},"size":501,"islog":false,"avg":1}
```

request cal
```
{"id":"rcal","t":0,"cmd":"rc","range":{"start":1000000,"end":4000000000},"size":3,"islog":false,"avg":1}
```
response
```
{"id":"rcal","t":0,"cmd":"rc","range":{"start":1000000,"end":4000000000},"size":3,"islog":false,"avg":1,"sparam":{"s11":false,"s12":false,"s21":false,"s22":false},"result":[{"s11":{"real":-0.00002794712781906128,"imag":-0.0002553611993789673},"s12":{"real":0.00004700571298599243,"imag":0.000007621943950653076},"s21":{"real":0.00005336105823516846,"imag":-0.000005081295967102051},"s22":{"real":0.000048279762268066406,"imag":-0.00021597743034362793},"freq":1000000},{"s11":{"real":-0.0000025406479835510254,"imag":-0.0002896711230278015},"s12":{"real":0.0015550479292869568,"imag":-0.00180024653673172},"s21":{"real":-0.0020975321531295776,"imag":0.001013830304145813},"s22":{"real":-0.0018929839134216309,"imag":0.0009096488356590271},"freq":2000500000},{"s11":{"real":0.002163603901863098,"imag":-0.008152559399604797},"s12":{"real":0.0015474259853363037,"imag":-0.0013009533286094666},"s21":{"real":0.0024291202425956726,"imag":0.0006962120532989502},"s22":{"real":0.0027683377265930176,"imag":-0.006796970963478088},"freq":4000000000}],"what":""}
```

request measurement
```
{"id":"dut1","t":0,"cmd":"crq","what":"dut1","avg":1,"sparam":{"s11":true,"s12":true,"s21":true,"s22":true}}
```
response
```
{"id":"dut1","t":0,"cmd":"crq","what":"dut1","avg":1,"sparam":{"s11":true,"s12":true,"s21":true,"s22":true},"result":[{"s11":{"real":0.7213804056329582,"imag":1.4104697588145223},"s12":{"real":-0.6894871206139953,"imag":0.2700840697023531},"s21":{"real":0.7031091027298517,"imag":0.47775336746611247},"s22":{"real":0.10341030075671172,"imag":0.3654223580872504},"freq":1000000},{"s11":{"real":-0.21755274620868625,"imag":-0.3754279747150785},"s12":{"real":0.6896210346136951,"imag":-0.12124783800081158},"s21":{"real":0.5913287788642378,"imag":-0.17444370058172678},"s22":{"real":-0.19807969221279273,"imag":-0.3255108569280837},"freq":2000500000},{"s11":{"real":-0.6383708479867005,"imag":-0.14863376649864046},"s12":{"real":0.29690952420750577,"imag":-0.0448141195835949},"s21":{"real":0.32891032123190567,"imag":0.003536763868706301},"s22":{"real":-0.4513970279816097,"imag":-0.14868940500662153},"freq":4000000000}]}
```


## Two port set up

## Pin connections

| Pin | Connection      |
|-----|-----------------|
| 2   | Switch Power    |
| 3   | Switch Power    |
| 4   | Port 2 Switch A |
| 5   | Port 2 Switch B |
| 6   | Port 2 Switch C |
| 8   | Port 1 Switch A |
| 9   | Port 1 Switch B |
| 10  | Port 1 Switch C |
| 15  | Port 2 LED      |
| 16  | Port 1 LED      |

## Switch Configuration

The states available in the prototype firmware were:

| State       | P1-A (D8) | P1-B (D9) | P1-C (D10) | P2-A (D4) | P2-B (D5) | P2-C (D6) |
|-------------|-----------|-----------|------------|-----------|-----------|-----------|
| P1 Short    | 0         | 0         | 1          | X         | X         | X         |
| P1 Open     | 0         | 1         | 0          | X         | X         | X         |
| P1 Load     | 0         | 1         | 1          | X         | X         | X         |
| P2 Short    | X         | X         | X          | 1         | 0         | 0         |
| P2 Open     | X         | X         | X          | 1         | 0         | 1         |
| P2 Load     | X         | X         | X          | 1         | 1         | 0         |
| Thru        | 1         | 0         | 0          | 0         | 1         | 1         |
| DUT-1 (5-2) | 1         | 0         | 1          | 0         | 1         | 0         |
| Ch0 P2      | X         | X         | X          | 0         | 0         | 0         |
| Ch1 P2      | X         | X         | X          | 0         | 0         | 1         |
| Ch7 P2      | X         | X         | X          | 1         | 1         | 1         |
| Ch0 P1      | 0         | 0         | 0          | X         | X         | X         |
| Ch6 P1      | 1         | 1         | 0          | X         | X         | X         |
| Ch7 P1      | 1         | 1         | 1          | X         | X         | X         |
| DUT-2 (6-1) | 1         | 1         | 0          | 0         | 0         | 1         |
| DUT-3 (7-0) | 1         | 1         | 1          | 0         | 0         | 0         |
| DUT-4 (0-7) | 0         | 0         | 0          | 1         | 1         | 1         |

The revised, minimal states are:

| State       | P1-A (D8) | P1-B (D9) | P1-C (D10) | P2-A (D4) | P2-B (D5) | P2-C (D6) |
|-------------|-----------|-----------|------------|-----------|-----------|-----------|
| P1 Short    | 0         | 0         | 1          | X         | X         | X         |
| P1 Open     | 0         | 1         | 0          | X         | X         | X         |
| P1 Load     | 0         | 1         | 1          | X         | X         | X         |
| P2 Short    | X         | X         | X          | 1         | 0         | 0         |
| P2 Open     | X         | X         | X          | 1         | 0         | 1         |
| P2 Load     | X         | X         | X          | 1         | 1         | 0         |
| Thru        | 1         | 0         | 0          | 0         | 1         | 1         |
| DUT-1       | 1         | 0         | 1          | 0         | 1         | 0         |
| DUT-2       | 1         | 1         | 0          | 0         | 0         | 1         |
| DUT-3       | 1         | 1         | 1          | 0         | 0         | 0         |
| DUT-4       | 0         | 0         | 0          | 1         | 1         | 1         |

Or in decimal

| State       | P1        | P2        | 
|-------------|-----------|-----------|
| Short       | 1         | 4         |
| Open        | 2         | 5         |
| Load        | 3         | 6         |
| Thru        | 4         | 3         |
| DUT-1       | 5         | 2         |
| DUT-2       | 6         | 1         |
| DUT-3       | 7         | 0         |
| DUT-4       | 0         | 7         |

This table implies the following connections:

### Port 1

| Channel   | Connection |
|-----------|------------|
| RF1       | DUT-4      |
| RF2       | Short      |
| RF3       | Open       |
| RF4       | Load       |
| RF5       | Thru       |
| RF6       | DUT-1      |
| RF7       | DUT-2      |
| RF8       | DUT-3      |

### Port 2

| Channel   | Connection |
|-----------|------------|
| RF1       | DUT-3      |
| RF2       | DUT-2      |
| RF3       | DUT-1      |
| RF4       | Thru       |
| RF5       | Short      |
| RF6       | Open       |
| RF7       | Load       |
| RF8       | DUT-4      |

These connections tally with the prototype.

Note that the connection on the HM321 RF chip on the datasheet.


### Current DUT Setup

The prototype as at 2023-03-02 is set up as follows:


| DUT | Type           | Port/Switch 1 | Port/Switch 2 | 
|-----|----------------|---------------|---------------|
| 1   | 2-Port Filter  | P1            | P2            |
| 2   | 4-Port Coupler | P1            | P4            | 
| 3   | 4-Port Coupler | P1            | P2            |
| 4   | 4-Port Coupler | P1            | P3            |

The 4-Port couplers are all the same, but with different port connections, to allow all four ports to be measured (albeit not as completely as a four-port VNA).

![hm321](./img/hm321.png)

## Current results

Set debug=true to get a table showing the outputs

```
        P1      P2     P1  P2
short: 1 0 0   0 0 1   1   4
open : 0 1 0   1 0 1   2   5
load : 1 1 0   0 1 1   3   6
thru : 0 0 1   1 1 0   4   3
dut1 : 1 0 1   0 1 0   5   2
dut2 : 0 1 1   1 0 0   6   1
dut3 : 1 1 1   0 0 0   7   0
dut4 : 0 0 0   1 1 1   0   7
```

## Developers

This section contains notes for developers. The code was refactored during 2023 to account for a bug that required two measurements to be taken in a row before the correct DUT was measured. 


### Measurements

A typical user request is for a calibrated measurement of a particular device under test (DUT).

 The `pkg/middle` receives that request over a websocket connection implemented by `pkg/stream`. Upon startup, the firmware makes a local connection to a `relay host` that is in turn connected to the `relay server` instance running in the cloud, that is turn connected to the user.

`pkg/middle` must take a number of measurements, both of cal standards and the DUT, then request a calibration calculation. Measurements require the RF switch to be set to connect the VNA to the right DUT or cal standard.

The `pkg/measure` package combines an RF switch (`pkg/rfusb`) and VNA instance (`pkg/pocket`) together so that it is straightforward to request a particular measurement. 

Once all the required measurements are taken (SOLT + DUT), then using the method shown in `pkg/calibrate`, the gRPC client in `pkjg/pb` is used to request the calibration from our local calibration server - see `py/server.py`. This needs to run locally to the firware because gRPC is HTTP/2 and that is not proxied by the cloud frontends available to us at present. The code in `pkg/pb` is autogenerated protocol buffer code.


```mermaid
 erDiagram
    USER }o--o{ CALIBRATED_MEASUREMENT : requests
    CALIBRATED_MEASUREMENT {
        string port
        array frequencyRange
        result  CalibratedSParametersDUT
    }
    CALIBRATED_MEASUREMENT ||--|| CALIBRATION : requires
        CALIBRATION {
            SParameters short
            SParameters open
            SParameters load
            SParameters thru
            SParameters DUT

        }
     CALIBRATION ||--|{ MEASUREMENT : requires
        MEASUREMENT {
                 string port
        }
    MEASUREMENT ||--|| HARDWARE : requires
        HARDWARE{
        pointer switch
            pointer VNA
        }    
```

TODO 

0. `middle` update hardware test to cycle through all dut at least twice, and then a third time in a different order
1. `middle` Use new direct switch connection in middle
2. `middle` Use new grpc connection in middle


### Golang files

The golang packages are in `./pkg`

If adding new subpackages, then update the reference to the project to the current branch, e.g. after adding the package `pb` to the `fix` branch, `go mod` couldn't find it `@latest` so it needed to be "upgraded" to look at the fix branch. Warning - after pulling to main, update to point to @latest again
```
go get -v -u github.com/practable/pocket-vna-two-port/pkg/pb@fix
go mod tidy
```
### Architecture

The RF Switch is connected directly rather than via `socat`+`websocat` as in the previous generation. This simplifies the configuration of the single board computer and avoids switch commands getting out of sequence on a channel.

The calibration is via gRPC call, again to avoid responses getting out of sequence over a channel. Note that gRPC uses HTTP/2 so we are probably stuck with running this locally on a container


### Building



### Testing


## Overview from the one-port repo

These notes should aid in understanding the ansible scripts in `./sbc`, even though some details differ. Note that the below information has NOT been edited to match changes in the current repo, so differs in the detail e.g. of calibration commands (now two-port), and S-parameter settings in the measurement commands. See commands above for the current version.

This repo contains code and design files for the key parts of the system.

- golang wrapper for the pocket-vna
- golang middleware and clients for the rf-switch and calibration services 
- python calibration service using scikit-rf & websocket libraries
- docker containerisation of the python calibration service
- C/C++ firmware for an arduino nano microcontroller
-
The largest part of this repo is the golang executable `vna` which sequences tasks between the four major parts: 
- user (via websocket service, routed through session host from the external practable.io ecosystem)
- vector network analyser (via golang-wrapped shared C library; local USB connection to hardware)
- rf switch controller (via websocket service, routed through session host; local USB connection to hardware)
- calibration service (via websocket service, routed through session host; runs on local docker container)

`vna` also supplies a `unlock` command to free up any pocketVNA connected to the computer. It is assumed that you only have one VNA per computer (it's possible to support more, and to choose explicitly, but that is not implemented. Our use case is one VNA per host computer,  and we want to auto-discover the VNA).

![system](./img/system.png)

For example, in order to perform a calibration, the golang code does the following tasks

- set the RF switch to the short cal standard
- scan the short, and store the data
- set the RF switch to the open cal standard
- scan the open, and store the data
- set the RF switch to the load cal standard
- scan the load, and store the data
- check all the data are of the correct length

To supply calibrated scans, the golang code does the following

- set the RF switch to the device under test (or whichever cal standard the user is interested in scanning)
- scan the device under test with the same parameters for which we have cal-standard data
- combine the dut data with the cal standard data, and send to the calibration service
- receive the results, re-format the data, and send it to the user

There is a <100ms performance overhead in sending the cal data to the calibration service each time, but this is much smaller than the scan times (tens of seconds), and avoids duplicating state (such as knowledge of what frequency points are being used, and when the cal standards were last scanned).

## Future expansion

Our preferred languages for each of the different parts of the system spanned golang, C/C++ and python. Instead of trying to tightly integrate them or duplicate functionality in one or the other, the JSON-over-websockets communication method was used to separate out concerns, facilitate testing, and allow replacing the RF switch circuit to support two-port cal, add extra calibration methods, and upgrade the API offered to the user to support the extra hardware.

## Installation 

This can be a bit finicky because the pocketVNA library has different name and requirements for different targets, and is only available as compiled source so cannot be automatically built for your target. The manufacturers support a number of major operating systems and platforms, so you should be ok. This code has been tested on amd64/linux and arm/linux.

### vna install 

The core of this repo is the `vna` executable.

0. Download this repo (`git clone git@github.com:timdrysdale/go-pocketvna.git`)
0. replace pkg/pocket/libPocketAPI.so with the appropriate pocketVNA library for your target
0. copy that library to your system path, e.g. for linux /usr/lib (for linux amd64/arm, add `.0` so the ending is `...so.0`)
0. get all dependencies `go get -u ./...' (you may need to add or remove the -lm flag for linking the maths library - not needed on linux 64 bit but needed on Raspbian buster
0. Compile the vna command (`go mod init github.com/timdrysdale/go-pocketvna && cd cmd/vna && go build`)
0. copy the `vna` executable to your usual location on the path (`go install` or do by hand)
0. check it is installed by trying `vna stream --help`

```
vna stream --help
Stream connects the first available pocketVNA to a websocket server. The websocket server is specified via an environment variable

export VNA_DESTINATION=ws://localhost:8888/ws/vna
export VNA_CALIBRATION=ws://localhost:8888/ws/calibration
export VNA_RFSWITCH=ws://localhost:8888/ws/rfswitch

vna stream 

Note that development can be enabled by setting environment variable VNA_DEVELOPMENT
export VNA_DEVELOPMENT=true

Usage:
  vna stream [flags]

Flags:
  -h, --help   help for stream
```


### session install

Install `session` from [practable.io](https://github.com/practable/relay)

If you are using the experiment remotely, you should configure a script to set the forwarding rules at start up (your practable.io admin contact will have details of the tokens etc your experiment needs)

session-rules
```
#!/bin/sh
dataTokenFile="/etc/practable/data.token"
dataAccessFile="/etc/practable/data.access"

dataToken=$(cat "$dataTokenFile")
dataAccess=$(cat "$dataAccessFile")

curl -X POST -H "Content-Type: application/json" -d '{"stream":"data","destination":"'"${dataAccess}"'","id":"1","token":"'"${dataToken}"'"}' http://localhost:8888/api/destinations 
```
This script can be run at start up by enabling a systemd service. The service file is in [services](./ansible/services)

### Docker install

Python is often used for system tasks, so we run our calibration service in a docker container to avoid any version and dependency conflicts.

The steps to manually install Docker are:

```
sudo apt update
# optional:
# sudo apt-upgrade && sudo reboot
curl -fsSL https://get.docker.com -o get-docker.sh
chmod +x ./get-docker.sh
sudo ./get-docker.sh

# for non-privileged use
dockerd-rootless-setuptool.sh install
#check you have docker running
docker version 
#sudo usermod -aG docker $(whoami) #this not be needed -
sudo chmod 666 /var/run/docker.sock #avoid permission denied
#sudo reboot #let all the changes take effect
#<wait for reboot>
```

Build and start the docker image (use tmux if you are ssh'ing into your rpi, so you can have extra terminals open)
```
export platform=linux-amd64  #or linux-arm32v7
cd py/docker/app-${platform}
./build.sh && ./start.sh
```

You can check the calibration service is up and running by connecting to it with `websocat`

```
$ websocat ws://localhost:8888/ws/calibration -
```

Calibration commands are usually quite long, so here is one with just a single data point
```
{"cmd": "oneport", "freq": [1000000.0], "short": {"real": [0.9166423490437918], "imag": [0.65760561459272446]}, "open": {"real": [0.8574903206586918], "imag": [0.43502949254752743]}, "load": {"real": [0.3002840906307519], "imag": [0.297151596182326]}, "dut": {"real": [0.4975782258013943], "imag": [0.4293572766329692]}}
```
which will return a response like this:


```
{"freq": [1000000.0], "S11": {"Real": 0.032134147957021554, "Imag": 0.0984021118681623}}
```

## vna setup
 
The usage / setup is mostly the same for a local system running on a desktop for testing, or deployed on a raspberry pi for remote access. The remote system tends to have everything autostart on boot, and include some setup to make a connection to the `practable.io` infrastructure servers. This is added complication that you do not need, and may not want, if working on a desktop. Starting the services manually is fine. Start `session host` first, then `websocat`, then you can start `socat`, `vna stream` and the `docker` calibration service in any order. The system will only supply calibrated data when the rfswitch and calibration services are running.

It is recommended to try unlocking the devices first - just in case the service is restarting and the pocketVNA was locked by the previous instance

The websocket connections between services are specified by environment variables, so that `vna stream` knows where to find them. The `session host` does not mind which particular endpoints are used - but the endpoints must be set the same for `vna stream` and the service that it is connecting to. The full instantiation of the [practable.io](https://practable.io) services on the host computer (usually a raspberry pi) are out of the scope of this documentation (early adopters are encouraged to contact us for information on how to obtain support <support@practable.io>).

Assuming your target is already running [relay/sessionhost](https://github.com/timdrysdale/relay) (see `cmd/session`), then you can use the following files to automate the startup of the services you need:

vna-data:
```
#!/bin/bash
export VNA_DESTINATION=ws://localhost:8888/ws/data
export VNA_CALIBRATION=ws://localhost:8888/ws/calibration
export VNA_RFSWITCH=ws://localhost:8888/ws/rfswitch
vna unlock
vna stream
```

socat-rfswitch:
```
#!/bin/sh
socat /dev/ttyUSB0,echo=0,b57600,crnl tcp:127.0.0.1:9999
```

Note you may have to add yourself to the dialout and tty groups.
```
sudo usermod -a -G dialout pi
sudo usermod -a -G tty pi
```

websocat-rfswitch:
```
#!/bin/sh
websocat ws://localhost:8888/ws/rfswitch tcp-listen:127.0.0.1:9999 --text
```

You can check the `rfswitch` is running by 
```
websocat ws://localhost:8888/ws/rfswitch -
```
Issue commands to change the port
```
{"set":"port","to":"open"}
```

should yield a reply of
```
{"report":"port","is":"open"}
```

### Service files

Service files are in <./ansible/services>. Note most are fairly similar, except `socat-data` which must wait long enough that `websocat-data` has started (ca 10s leaves a safe margin in practice).

### Verifying that the services are working

Connect to `vna` over websocket using `websocat`
```
websocat ws://localhost:8888/ws/data -
```
Then send a ReasonableFrequencyRange command
```
{"cmd":"rr"}
```
which should return something like
```
{"id":"","t":0,"cmd":"rr","range":{"start":500000,"end":4000000000}}
```

This confirms that `vna` is up and running, and connected to a pocketVNA.

Now try a command that needs the rfswitch
```
{"cmd":"rc","range":{"start":1000000,"end":4000000000},"size":3,"islog":false,"avg":1,"sparam":{"s11":true,"s12":false,"s21":false,"s22":false}}
```
should return something like:

```
{"id":"","t":0,"cmd":"rc","range":{"start":1000000,"end":4000000000},"size":3,"islog":false,"avg":1,"sparam":{"s11":true,"s12":false,"s21":false,"s22":false},"result":[{"s11":{"real":-0.0007592514157295227,"imag":0.0006944984197616577},"s12":{"real":0,"imag":0},"s21":{"real":0,"imag":0},"s22":{"real":0,"imag":0},"freq":1000000},{"s11":{"real":0.0004634261131286621,"imag":-0.0001447424292564392},"s12":{"real":0,"imag":0},"s21":{"real":0,"imag":0},"s22":{"real":0,"imag":0},"freq":2000500000},{"s11":{"real":0.0024859830737113953,"imag":-0.005337625741958618},"s12":{"real":0,"imag":0},"s21":{"real":0,"imag":0},"s22":{"real":0,"imag":0},"freq":4000000000}]}

```

and a command that needs the calibration service

```
{"cmd":"crq","what":"dut","avg":1,"sparam":{"s11":true,"s12":false,"s21":false,"s22":false}}
```

which should return something like:

```
{"id":"","t":0,"cmd":"crq","what":"dut","avg":1,"sparam":{"s11":true,"s12":false,"s21":false,"s22":false},"result":[{"s11":{"real":1.0031062693348898,"imag":-0.012728241016514735},"s12":{"real":0,"imag":0},"s21":{"real":0,"imag":0},"s22":{"real":0,"imag":0},"freq":1000000},{"s11":{"real":-0.1335237960932827,"imag":-0.09618931729902436},"s12":{"real":0,"imag":0},"s21":{"real":0,"imag":0},"s22":{"real":0,"imag":0},"freq":2000500000},{"s11":{"real":-0.4883905666762738,"imag":-0.13897086282584725},"s12":{"real":0,"imag":0},"s21":{"real":0,"imag":0},"s22":{"real":0,"imag":0},"freq":4000000000}]}
```

## Interface 

There are three commands directly supported by the pocketvna shared library:

0. `rr`: get the reasonable frequency range for the device (useful `proof of connection` command when troubleshooting)
0. `sq`: get S-parameters at a single frequency (uncalibrated)
0. `rq`: get S-paramters at a range of frequencies (uncalibrated)

we add two more commands to support calibration, which are the two main commands we use: 

0. `rc`: do a calibration over a frequency range for a certain number of steps 

0. `crq`: get S-parameters frequencies (calibrated), using whatever frequency points were specified in the last `rc`

There is also a heartbeat sent every second from the driver, to let you know it is still connected. You can use the absence of this heartbeat to infer that your connection has dropped (it is unlikely the driver has stopped, although that is technically possible if there is a power outage). This heartbeat is added because VNA experiments are unlikely to use a camera (and we've been using the video to check for connection drops, so without it we need something else to check).

```
{"cmd":"hb"}
```

In order to relate commands to responses you can include an ID `id` (string) and/or time `t` (int) field.

### rr



Command

```
{"cmd":"rr"}
```

Response

```
{"id":"","t":0,"cmd":"rr","range":{"Start":500000,"End":4000000000}}
```

Start is the lowest frequency (in Hz) that the VNA can operate at, and End is the highest.

### sq

`sq` gets the requested S-parameters at a single frequency. 

The parameters are:

    - freq (in Hz) at which to take the measurement at
	- avg number of times to take the reading (which are then averaged together)
	- sparam select which S-params are needed - the more you select, the longer it takes

Command:

```
{"id":"945102d5-94e4-448e-bbbf-48384c662711","t":1634664795,"cmd":"sq","freq":100000,"avg":1,"sparam":{"S11":true,"S12":false,"S21":true,"S22":false}}
```
Response:

```
{"id":"945102d5-94e4-448e-bbbf-48384c662711","t":1634664795,"cmd":"sq","freq":100000,"avg":1,"sparam":{"S11":true,"S12":false,"S21":true,"S22":false},"result":{"S11":{"Real":0.0001702234148979187,"Imag":0.0005754455924034119},"S12":{"Real":0,"Imag":0},"S21":{"Real":-0.00004191696643829346,"Imag":-0.00012067705392837524},"S22":{"Real":0,"Imag":0}}}
```

### rq

Get the requested S-parameters at a range of frequencies. Parameters:

    - range (start, end) the start and end of the frequency range to take measurements
	- size the number of separate frequency points to measure at, across the range
	- isLog whether to use a log (isLog = true) or linear (isLog = false) distribution
	- avg number of times to take each reading (which are then averaged together)
	- sparam select which S-params are needed - the more you select, the longer it takes

command
```
{"cmd":"rq","range":{"Start":100000,"End":4000000},"size":2,"isLog":true,"avg":1,"sparam":{"S11":true,"S12":false,"S21":true,"S22":false}}
```

response
```
{"id":"","t":0,"cmd":"rq","range":{"Start":100000,"End":4000000},"size":2,"isLog":true,"avg":1,"sparam":{"S11":true,"S12":false,"S21":true,"S22":false},"result":[{"S11":{"Real":0.00013846158981323242,"Imag":0.00027057528495788574},"S12":{"Real":0,"Imag":0},"S21":{"Real":-0.000031754374504089355,"Imag":-0.0002350062131881714},"S22":{"Real":0,"Imag":0}},{"S11":{"Real":0.00470772385597229,"Imag":0.003948085010051727},"S12":{"Real":0,"Imag":0},"S21":{"Real":0.000017777085304260254,"Imag":-0.000005081295967102051},"S22":{"Real":0,"Imag":0}}]}
```

Note that the frequencies are not (currently) included in the response. They can be calculated according to the formulas given in the openAPI library header, which are (@p denotes a parameter, e.g. ```@p start``` means the parameter start)

```
     *  Distributions:
     *    Linear:      (@p start * 1000 + ((@p end - @p start) * 1000 / (@p steps - 1)) * index) / 1000
     *       (Pay Attention: all numbers are integers. Last element is forced to be equalt to @p end)
     *    Logarithmic: (@p from * powf((float)to / from, (float)index / (steps - 1)))
     *       Formula is taken from 10 ** (lg from +  ( lg to - lg from ) * index /  (steps - 1)). 4-bytes float are used
     *       Pay attention: arithmetic is pretty imprecise on a device
	 
```
These frequency points are calculated in this library, because they are not sent from the hardware. Here is a table of 11 points from 1MHz - 500MHz, taken from the native application for Linear and Log scans. Note that these values from the native app are not quite the same as we expect from the formula for the case of a Log distribution, although the precision issue is acknowledged in the library header file from the supplier as if it is a hardware issue but the values are probably calculated in the native app - the differences are quite small and probably not significant.


| Point  | Linear (both)     |  Log (native app)      | Log (go-pocketVNA)      |  delta |
| -------|-------------|------------|-----------|--------|
| 0      |  1000000    | 1000000    | 1000000   |        |
| 1      |  50900000   | 1861646    | 1861646   |        |
| 2      |  100800000  | 3465724    | 3465724   |        |
| 3      |  150700000  | 6451951    | 6451950   |  -1    |
| 4      |  200600000  | 12011245   | 12011244  |  -1    |
| 5      |  250500000  | 22360680   | 22360680  |        |
| 6      |  300400000  | 41627668   | 41627660  |  -8    |
| 7      |  350300000  | 77495944   | 77495949  |  +5    |
| 8      |  400200000  | 144270000  | 144269991 |  -9    |
| 9      |  450100000  | 268579552  | 268579588 | +36    |
| 11     |  500000000  | 500000000  | 500000000 |        |


### rc

This uses the same JSON structure as the `rq` command, but change `"cmd":"rq"` to `"cmd":"rc"` and ensure that `"sparam":{"S11":true,"S12":false,"S21":false,"S22":false}}` because only cal on port one is supported:

```
{"cmd":"rc","range":{"Start":100000,"End":4000000},"size":2,"isLog":true,"avg":1,"sparam":{"S11":true,"S12":false,"S21":false,"S22":false}}
```

The maximum size is 512, and it is conventional to ask for 501 points for a nice even spacing. A calibration for 501 points takes approx 30 seconds.

### crq

This is a subset of the `rq` object, with `"cmd":"crq"`, and ensure that `"sparam":{"S11":true,"S12":false,"S21":false,"S22":false}}`

```
{"cmd":"crq","avg":1,"sparam":{"S11":true,"S12":false,"S21":false,"S22":false}}
```


## Example data

Good data from a passive device under test always returns a gain of less than 0dB (i.e. without power from a bias supply, it can't add any power of its own) so we can see that in the following results, we should not rely on values close to DC (for the calibration we did, at least). Also, there is some drift in this measurement of a cable, bearing in mind the cable was re-connected after the first few scans and the device moved, so not an entirely fair test on the device. This data was collected in the native app and plotted using python's matplotlib.pyplot.

![magnitude](./img/cable-dB.png)
![phase](./img/cable-deg.png)

### Raw data 

```
{"cmd":"rq","range":{"Start":100000,"End":4000000},"size":201,"isLog":true,"avg":1,"sparam":{"S11":true,"S12":true,"S21":true,"S22":true}}
```

Result returns about 15 sec later, with S-parameters at 201 frequency points:

```
{"id":"","t":0,"cmd":"rq","range":{"Start":100000,"End":4000000},"size":201,"isLog":true,"avg":1,"sparam":{"S11":true,"S12":true,"S21":true,"S22":true},"result":[{"S11":{"Real":0.000034302473068237305,"Imag":0.00037600845098495483},"S12":{"Real":-0.0010073482990264893,"Imag":0.00012448430061340332},"S21":{"Real":-0.0013414323329925537,"Imag":0.0005551204085350037},"S22":{"Real":0.00003937631845474243,"Imag":0.00023628026247024536},"Freq":100000},{"S11":{"Real":0.00006987154483795166,"Imag":0.00010670721530914307},"S12":{"Real":-0.0010479986667633057,"Imag":0.00009654462337493896},"S21":{"Real":-0.0013172999024391174,"Imag":0.0008002892136573792},"S22":{"Real":0.00007876008749008179,"Imag":0.0003162994980812073},"Freq":101862},{"S11":{"Real":0.00014227628707885742,"Imag":0.00035568326711654663},"S12":{"Real":-0.0010594278573989868,"Imag":0.0003429800271987915},"S21":{"Real":-0.0013439729809761047,"Imag":0.0009235069155693054},"S22":{"Real":0.00011432915925979614,"Imag":0.0001664087176322937},"Freq":103758}, ... <snip> ...]}
```



## Developer info

These sections contain some information that may be relevant to developers.

### CGO info

An example of using cgo can be found [here](https://github.com/timhughes/cgoexample) - note you have to delete hello.c from the project directory to avoid `multiple definition of main` error.

There are some differences between C and golang which could create difficulties for users of some C libraries [here](https://docs.yottadb.com/Presentations/DragonsofCGO.pdf), especially where there are:
0. pointers to structs containing other pointers
1. callbacks
2. issues not apparent until garbage collection is triggered

### Performance

It is unclear whether the newer, faster firmware is supported on linux. From the manual:-

```
Pay attention : that newer firmware supports 2 interfaces:
1. HID -- it is universal and easy. Supported by all platforms. But it is slow
2. VCI -- available since firmware V2.10. On Windows requires a special driver (for
example zadig/Interface 1) It is faster. May be unavailable on Mac-OS
In this case both interfaces are listed as individual (separated\independent) items with
the same Serial Number.
It is bad idea to connect both interfaces at the same time and especially bad to perform
scan concurrently
```

Thus to use this newer firmware we will likely need to enumerate devices and handle lists, etc. This should be ok because the types in the C library can be used as parameters.

### Descriptions

Note that getting the first handle will literally just get a handle - you cannot access the description from that handle. If you want the description, you need to get a list of descriptions, and then hask for the handle for the description you like the best. So don't bother trying to access description fields from the handle or handle pointer.

#### Result codes

The result strings are copied from `pocket.h`, primarily to help with debugging by providing some meaning to non-zero error codes. The result from the functions is an C enum, which is only an `int`, so can be typecast back to Go int, then used as the index into an array of strings. For example, If there is no device, then attempting to get a device handle returns code 0x05, or `PVNA_Res_NoDevice` or if a function completes successfully it returns 0x00, which has the string `PVNA_Res_Ok`.

### Testing

The testing for this library assumes the presence of pocketVNA - this avoids mocking the API library and the hardware.

The testing routine issues are force unlock command which should unlock any available pocketVNA. Testing will take place with the first available handle found by the openAPI library. If there is no hardware available you will see:
```
=== RUN   TestGetReleaseHandle
    TestGetReleaseHandle: pocket_test.go:45: 
        	Error Trace:	pocket_test.go:45
        	Error:      	Received unexpected error:
        	            	PVNA_Res_NoDevice
        	Test:       	TestGetReleaseHandle
--- FAIL: TestGetReleaseHandle (0.00s)
```

This probably just means your pocketVNA is unplugged.

If you plug it in and it is working, then you will see the test(s) pass (this output before other tests were added)

```
go build -v && go test -v && go vet
=== RUN   TestGetReleaseHandle
--- PASS: TestGetReleaseHandle (0.06s)
PASS
ok  	_/home/your-path/go-pocketvna/pkg	0.100s
```

#### SingleQuery

Cable disconnected at port 1

```
=== RUN   TestSingleQuery
{(-3.302842378616333e-05-0.00019435584545135498i) (1.2695789337158203e-05-1.9058585166931152e-05i) (-7.614493370056152e-06-2.4139881134033203e-05i) (3.810226917266846e-05-0.00016132742166519165i)}
```

Cable connected

```
{(-0.0013300031423568726+0.000985749065876007i) (0.0026688948273658752-0.005279362201690674i) (0.002047717571258545-0.006369277834892273i) (-0.0012982413172721863+0.0008942857384681702i)}
```

### How to load a touchstone file in python

Install [scikit-rf](https://scikit-rf.readthedocs.io/en/latest/tutorials/Installation.html), usually
```
pip install scikit-rf

```

start a python session, load and plot your file (change to suit your filename)

```
import matplotlib.pyplot as plt
import skrf as rf
load = rf.Network('dut.s1p')
dut.plot_s_db()
plt.savefig('dut_db.png', dpi=300)
plt.close()
dut.plot_s_deg()
plt.savefig('dut_deg.png', dpi=300)
plt.close()
```


### API to the API

since pointers container pointers are difficult to handle, the preferred API is the range API, probably.

## Testability, and future substitution of other hardware ...

There is really only three commands that the user interface needs to know about - GetReasonableFrequencyRange, SingleQuery, RangeQuery.

Therefore a simple remote driver would grab a device handle, and sit back and wait for these commands.

Getting the callback to work would be pretty useful, too, for the rangequery command. Can we poke percentage complete values into a channel?


### Small gotchas with cgo

#### Could not determine type of name error

```
To access a struct, union, or enum type directly, prefix it with struct_, union_, or enum_, as in C.struct_stat.
```
for example
```
C.enum_PocketVNADistribution(distr), //note we have to add enum_ to access this name
```


### Preparing PI image

Use [pishrink](https://github.com/Drewsif/PiShrink)

```
sudo dd if=/dev/sdb of=/home/<user>/rpi-config/vna-cal-2022-01-31.img bs=4M status=progress
sudo pishrink.sh -p vna-cal-2022-01-31.img vna-cal-2022-01-31-shrunk.img
```

