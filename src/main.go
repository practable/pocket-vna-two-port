//
// main.go
// Copyright (C) 2019 Tim Hughes
//
// Distributed under terms of the MIT license.
//
package main

/*
#cgo CFLAGS: -g -Wall
#cgo LDFLAGS: -L. -lPocketVnaApi_x64
#include "pocketvna.h"
*/
import "C"
import "fmt"

func main() {

	result := C.pocketvna_force_unlock_devices()
	fmt.Println(decode(int(result)))
	handle := C.PVNA_DeviceHandler(nil)
	result = C.pocketvna_get_first_device_handle(&handle)
	fmt.Println(decode(int(result)))
	result = C.pocketvna_release_handle(&handle)
	fmt.Println(decode(int(result)))
}

func decode(code int) string {

	names := [...]string{
		"PVNA_Res_Ok",
		"PVNA_Res_NoDevice",
		"PVNA_Res_NoMemoryError",
		"PVNA_Res_CanNotInitialize",
		"PVNA_Res_BadDescriptor",
		"PVNA_Res_DeviceLocked",
		"PVNA_Res_NoDevicePath",
		"PVNA_Res_NoAccess",
		"PVNA_Res_FailedToOpen",
		"PVNA_Res_InvalidHandle",
		"PVNA_Res_BadTransmission",
		"PVNA_Res_UnsupportedTransmission",
		"PVNA_Res_BadFrequency",
		"PVNA_Res_DataReadFailure",
		"PVNA_Res_EmptyResponse",
		"PVNA_Res_IncompleteResponse",
		"PVNA_Res_FailedToWriteRequest",
		"PVNA_Res_ArraySizeTooBig",
		"PVNA_Res_BadResponse",
		"PVNA_Res_DeviceResponseSection",
		"PVNA_Res_Response_UNKNOWN_MODE",
		"PVNA_Res_Response_UNKNOWN_PARAMETER",
		"PVNA_Res_Response_NOT_INITIALIZED",
		"PVNA_Res_Response_FREQ_TOO_LOW",
		"PVNA_Res_Response_FREQ_TOO_HIGH",
		"PVNA_Res_Response_OutOfBound",
		"PVNA_Res_Response_UNKNOWN_VARIABLE",
		"PVNA_Res_Response_UNKNOWN_ERROR",
		"PVNA_Res_Response_BAD_FORMAT",
		"PVNA_Res_ExtendedSection",
		"PVNA_Res_ScanCanceled",
		"PVNA_Res_Rfmath_Section",
		"PVNA_Res_No_Data",
		"PVNA_Res_LIBUSB_Error",
		"PVNA_Res_LIBUSB_CanNotSelectInterface",
		"PVNA_Res_LIBUSB_Timeout",
		"PVNA_Res_LIBUSB_Busy",
		"PVNA_Res_VCI_PrepareScanError",
		"PVNA_Res_VCI_Response_Error",
		"PVNA_Res_EndLEQStart",
		"PVNA_Res_VCI_Failed2OpenProbablyDriver",
		"PVNA_Res_HID_AdditionalError",
		"PVNA_Res_Fail",
	}

	if code == 255 {
		return names[len(names)-1]

	} else {
		return names[code]
	}

}
