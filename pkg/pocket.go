package pocket

/*
#cgo CFLAGS: -g -Wall
#cgo LDFLAGS: -L. -lPocketVnaApi_x64
#include "pocketvna.h"
*/
import "C"
import (
	"errors"
)

var Results = [...]string{
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

func ForceUnlockDevices() error {

	result := C.pocketvna_force_unlock_devices()

	return decode(result)

}

func GetFirstDeviceHandle() (C.PVNA_DeviceHandler, error) {

	handle := C.PVNA_DeviceHandler(nil)
	result := C.pocketvna_get_first_device_handle(&handle)
	return handle, decode(result)

}

func ReleaseHandle(handle C.PVNA_DeviceHandler) error {

	result := C.pocketvna_release_handle(&handle)
	return decode(result)

}

/* @brief Get reasonable frequency range IOW a range device can process correctly
   Usually it is narrower than [1_Hz; 6_GHz].

       @ingroup API
       @param handle  A pointer to Device.
       @param from    A pointer (reference) where to save lowest frequency a device can process correctly
       @param to      A pointer (reference) where to save highest frequency a device can process correctly

       @returns
           This function returns Result: 'Ok' on success, 'PVNA_Res_InvalidHandle' if handle is invalid

   PVNA_EXPORTED PVNA_Res   pocketvna_get_reasonable_frequency_range(const PVNA_DeviceHandler handle, PVNA_Frequency * from, PVNA_Frequency * to);
*/

func GetReasonableFrequencyRange(handle C.PVNA_DeviceHandler) (uint64, uint64, error) {

	from := C.PVNA_Frequency(0)
	to := C.PVNA_Frequency(0)
	result := C.pocketvna_get_reasonable_frequency_range(handle, &from, &to)

	return uint64(from), uint64(to), decode(result)

}

func decode(result C.PVNA_Res) error {

	code := int(result)

	if code == 0 {
		return nil
	} else {

		if code == 255 {
			return errors.New(Results[len(Results)-1])
		} else {
			return errors.New(Results[code])
		}
	}
}
