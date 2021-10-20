# go-pocketvna
A go-wrapper for the pocketvna.com 2-port vector network analyser's openAPI

## Background

The pocketVNA has a so-called 'open API' which is - contrary to open-source principles - only distributed in compiled form. Nonetheless, this compiled shared library written in pure-C is available for a number of platforms, a colleague already bought some pocketVNA, and we want to use them remotely, so use this library we must. Fortunately, Cgo allows Go code to call C-libraries, so I can write websocket and JSON messaging code in golang rather than C.


## Installation 

This is a bit finicky because the pocketVNA library has different name and requirements for different targets, and is only available as compiled source so cannot be automatically built for your target.

0. Download the repo
0. replace pkg/pocket/libPocketAPI.so with the appropriate pocketVNA library for your target
0. copy that library to your system path, e.g. for linux /usr/lib (for linux, add `.0`)
0. get all dependencies `go get -u ./...' (you may need to add or remove the -lm flag for linking the maths library - not needed on linux 64 bit but needed on Raspbian buster
0. Compile the vna command (cd to `cmd/vna` then `go build`)
0. copy the `vna` executable to your usual location on the path

## Usage

The streaming target is set by an environment variable. It is recommended to try unlocking the devices first - just in case the service is restarting and the pocketVNA was locked by the previous instance.

Assuming your target is running [relay/sessionhost](https://github.com/timdrysdale/relay) (see `cmd/session`), then you can use the following files to set up streaming that resumes on each boot:


vna-data:
```
#!/bin/bash
export VNA_DESTINATION=ws://localhost:8888/ws/data
vna unlock
vna stream
```

vna-data.service:
```
[Unit]
Description=get data from pocket vna 
After=network.target session.service 
Wants=session.service 

[Service]
Restart=on-failure
RestartSec=5s
ExecStartPre=/bin/sleep 1
ExecStart=/usr/local/bin/vna-data

[Install]
WantedBy=multi-user.target
```

```
#!/bin/sh
dataTokenFile="/etc/practable/data.token"
dataAccessFile="/etc/practable/data.access"

dataToken=$(cat "$dataTokenFile")
dataAccess=$(cat "$dataAccessFile")

curl -X POST -H "Content-Type: application/json" -d '{"stream":"data","destination":"'"${dataAccess}"'","id":"1","token":"'"${dataToken}"'"}' http://localhost:8888/api/destinations 
```

## Interface 

There are three commands

0. rr: get the reasonable frequency range for the device 
0. sq: get S-parameters at a single frequency 
0. rq: get S-paramters at a range of frequencies


In order to relate commands to responses you can include an id (string) and/or time (int) field.

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

Get the requested S-paramters at a range of frequencies. Parameters:

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


| Point  | Linear      |  Log       | Log       |  delta |
|        |  (both)     | native app | this lib  |        |
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
