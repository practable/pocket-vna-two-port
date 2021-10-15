/*
Package pocket uses cgo to wrap the shared C library for the pocketVNA openAPI

The commands supported are

ForceUnlock
GetFirstDeviceHandle
ReleaseHandle
GetReasonableFrequencyRange
SingleQuery
RangeQuery

Function call result codes are decoded as required, into strings as specified in pocket.h

*/

package pocket

/*
#cgo CFLAGS: -g -Wall
#cgo LDFLAGS: -L. -lPocketVnaApi_x64
#include "pocketvna.h"
*/
import "C"
import (
	"context"
	"errors"
)

// does not compile if in types.go ("C undefined")
type VNA struct {
	handle C.PVNA_DeviceHandler
}

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

func getFirstDeviceHandle() (C.PVNA_DeviceHandler, error) {

	handle := C.PVNA_DeviceHandler(nil)
	result := C.pocketvna_get_first_device_handle(&handle)
	return handle, decode(result)

}

func releaseHandle(handle C.PVNA_DeviceHandler) error {

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

func getReasonableFrequencyRange(handle C.PVNA_DeviceHandler) (uint64, uint64, error) {

	from := C.PVNA_Frequency(0)
	to := C.PVNA_Frequency(0)
	result := C.pocketvna_get_reasonable_frequency_range(handle, &from, &to)

	return uint64(from), uint64(to), decode(result)

}

/*  @brief Query device for some Network Parameters for particular frequency
     *
     *  It accepts @p handle and gets Network parameters @p params

        @ingroup API
        @param handle    A pointer to Device.
        @param frequency A frequency value. Usually it should be between [1_Hz; 6_GHz]
        @param average   A average times to ask hardware. Usually should be between [1; 1000]
        @param params    Network Parameters that should be taken: S11 or S21 or S12 or S22. Use '|' to combine
        @param s11       Pointer to SParam structure (pair of double). S11 Network Parameter will be here is @p params asked for it
        @param s21       Pointer to SParam structure (pair of double). S21 Network Parameter will be here is @p params asked for it
        @param s12       Pointer to SParam structure (pair of double). S21 Network Parameter will be here is @p params asked for it
        @param s22       Pointer to SParam structure (pair of double). S22 Network Parameter will be here is @p params asked for it

        @returns
            This function returns Result: 'Ok' on success, 'PVNA_Res_InvalidHandle' if handle is invalid, or any other 'Result'

    PVNA_EXPORTED PVNA_Res   pocketvna_single_query(const PVNA_DeviceHandler handle,
                                          const PVNA_Frequency frequency,
                                          const uint16_t average, const PVNA_NetworkParam params,
                                          PVNA_Sparam * s11,  PVNA_Sparam * s21,
                                          PVNA_Sparam * s12,  PVNA_Sparam * s22);
typedef struct ImitComplexD {
    double real;
    double imag;
} PVNA_Sparam;

enum PocketVnaTransmissionEnum{ PVNA_SNone = 0x00,
                                PVNA_S21   = 0x01,
                                PVNA_S11   = 0x02,
                                PVNA_S12   = 0x04,
                                PVNA_S22   = 0x08,

                                PVNA_FORWARD= PVNA_S11 | PVNA_S21,
                                PVNA_BACKWARD=PVNA_S12 | PVNA_S22,
                                PVNA_ALL   = PVNA_FORWARD | PVNA_BACKWARD
};

typedef enum PocketVnaTransmissionEnum PVNA_NetworkParam;
*/

func encodeParams(p SParamSelect) C.PVNA_NetworkParam {

	n := 0

	if p.S21 {
		n += 1
	}
	if p.S11 {
		n += 2
	}
	if p.S12 {
		n += 4
	}
	if p.S22 {
		n += 8
	}

	return C.PVNA_NetworkParam(n)

}

func singleQuery(handle C.PVNA_DeviceHandler, freq uint64, avg uint16, p SParamSelect) (SParam, error) {

	S11 := C.PVNA_Sparam{0.0, 0.0}
	S12 := C.PVNA_Sparam{0.0, 0.0}
	S21 := C.PVNA_Sparam{0.0, 0.0}
	S22 := C.PVNA_Sparam{0.0, 0.0}

	result := C.pocketvna_single_query(handle, C.PVNA_Frequency(freq), C.uint16_t(avg), encodeParams(p), &S11, &S21, &S12, &S22)

	s := SParam{
		S11: complex(S11.real, S11.imag),
		S12: complex(S12.real, S12.imag),
		S21: complex(S21.real, S21.imag),
		S22: complex(S22.real, S22.imag),
	}

	return s, decode(result)
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

/*   * @brief Query device for some Network Parameters using a distribution formula
     *
     *   It accepts @p handle and gets Network parameters @p params. Frequency point is calculated by distribution formula
     *  Distributions:
     *    Linear:      (@p start * 1000 + ((@p end - @p start) * 1000 / (@p steps - 1)) * index) / 1000
     *       (Pay Attention: all numbers are integers. Last element is forced to be equalt to @p end)
     *    Logarithmic: (@p from * powf((float)to / from, (float)index / (steps - 1)))
     *       Formula is taken from 10 ** (lg from +  ( lg to - lg from ) * index /  (steps - 1)). 4-bytes float are used
     *       Pay attention: arithmetic is pretty imprecise on a device

        @ingroup API
        @param handle   A pointer to Device
        @param start    Start Frequency
        @param end      End Frequency. Should be greater than @p start
        @param steps    Number of frequency points
        @param distr    A code for distribution formula (Linear)
        @param average  A average times to ask hardware. Usually should be between [1; 1000]
        @param params   Network Parameters that should be taken: S11 or S21 or S12 or S22. Use '|' to combine
        @param s11a      Array to SParam structures (pairs of double). S11 Network Parameters will be here is @p params asked for it
        @param s21a      Array to SParam structures (pairs of double). S21 Network Parameters will be here is @p params asked for it
        @param s12a      Array to SParam structures (pairs of double). S21 Network Parameters will be here is @p params asked for it
        @param s22a      Array to SParam structures (pairs of double). S22 Network Parameters will be here is @p params asked for it
        @param progress  Callback structure. It if is not NULL callee will be notified about currently processed index of frequency

        @returns

    PVNA_EXPORTED PVNA_Res   pocketvna_range_query(
            const PVNA_DeviceHandler handle,
            const PVNA_Frequency start, const PVNA_Frequency end, const uint32_t size, enum PocketVNADistribution distr,
            const uint16_t average, const PVNA_NetworkParam params,
            PVNA_Sparam * s11a, PVNA_Sparam * s21a,
            PVNA_Sparam * s12a, PVNA_Sparam * s22a,
            PVNA_ProgressCallBack * progress
    );

enum PocketVNADistribution {
    PVNADist_Linear=1,
    PVNADist_Log=2
};

*/
type Distribution int

const (
	Undefined Distribution = iota //handle default value being undefined
	Linear
	Log
)

// We do not implement the callback for this version ...
func rangeQuery(handle C.PVNA_DeviceHandler, start, end uint64, size int, distr int, avg uint16, p SParamSelect) ([]SParam, error) {

	S11 := [512]C.PVNA_Sparam{}
	S12 := [512]C.PVNA_Sparam{}
	S21 := [512]C.PVNA_Sparam{}
	S22 := [512]C.PVNA_Sparam{}

	result := C.pocketvna_range_query(handle,
		C.PVNA_Frequency(start),
		C.PVNA_Frequency(end),
		C.uint32_t(size),
		C.enum_PocketVNADistribution(distr), //note we have to add enum_ to access this name
		C.uint16_t(avg),
		encodeParams(p),
		&S11[0],
		&S21[0],
		&S12[0],
		&S22[0],

		nil)

	ss := []SParam{}

	for i := 0; i < int(size); i++ {

		s := SParam{
			S11: complex(S11[i].real, S11[i].imag),
			S12: complex(S12[i].real, S12[i].imag),
			S21: complex(S21[i].real, S21[i].imag),
			S22: complex(S22[i].real, S22[i].imag),
		}

		ss = append(ss, s)

	}

	return ss, decode(result)

}

func (v *VNA) RangeQuery(r RangeQuery) (RangeQuery, error) {

	distr := 0

	if r.LogDistribution {
		distr = 1
	}

	sparams, err := rangeQuery(v.handle, r.Range.Start, r.Range.End, r.Size, distr, r.Avg, r.Select)

	if err != nil {
		return r, err
	}

	r.Result = sparams

	return r, err

}

/* RunHandle provides a go channel interface to a given instance of a pocket VNA device

There are two uni-directional channels, one to receive commands, the other to reply with data

*/

func (c *VNA) Run(command <-chan interface{}, result chan<- interface{}, e chan<- error, ctx context.Context) {

}

/*
func RunFirstAvailable(c <- chan Command, d chan <- Data, e chan <- errors.Error, ctx context.Context) {

	handle, err := GetFirstDeviceHandle()

	from, _, err := GetReasonableFrequencyRange(handle)

	assert.NoError(t, err)

	s, err := SingleQuery(handle, from, 1, SParamSelect{true, true, true, true})

	assert.NoError(t, err)

	fmt.Println(s)

	err = ReleaseHandle(handle)
	assert.NoError(t, err)

	RunHandle(h C.PVNA_DeviceHandler, c chan<- Command, d <-chan Data, ctx context.Context)

}
*/
