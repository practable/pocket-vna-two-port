# go-pocketvna
A go-wrapper for the pocketvna.com 2-port vector network analyser's openAPI

## Background

The pocketVNA has a so-called 'open API' which is - contrary to open-source principles - only distributed in compiled form. Nonetheless, this compiled shared library written in pure-C is available for a number of platforms, a colleague already bought some pocketVNA, and we want to use them remotely, so use this library we must. Fortunately, Cgo allows Go code to call C-libraries, so I can write websocket and JSON messaging code in golang rather than C.


### Result codes

The result strings are copied from `pocket.h`, primarily to help with debugging by providing some meaning to non-zero error codes. The result from the functions is an C enum, which is only an `int`, so can be typecast back to Go int, then used as the index into an array of strings. For example, If there is no device, then attempting to get a device handle returns code 0x05, or `PVNA_Res_NoDevice` or if a function completes successfully it returns 0x00, which has the string `PVNA_Res_Ok`.





