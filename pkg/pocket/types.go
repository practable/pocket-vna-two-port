package pocket

import (
	"context"
	"time"
)

type VNAService struct {
	VNA      *VNA
	Ctx      context.Context
	Request  chan interface{}
	Response chan interface{}
	Timeout  time.Duration
}

/*
typedef struct PocketVnaDeviceDesc {
    const char * path;
    PVNA_Access access;

    const wchar_t * serial_number;

    const wchar_t * manufacturer_string;
    const wchar_t * product_string;

    uint16_t release_number;

    uint16_t pid;
    uint16_t vid;
    uint16_t ciface_code; //value from ConnectionInterfaceCode

    struct PocketVnaDeviceDesc * next;
} PVNA_DeviceDesc;
*/
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

type Description struct {
	Serial       string `json:"serial"`
	Manufacturer string `json:"manufacturer"`
	Product      string `json:"product"`
	Release      int    `json:"release"`
}

type SParamSelect struct {
	S11 bool `json:"s11"`
	S12 bool `json:"s12"`
	S21 bool `json:"s21"`
	S22 bool `json:"s22"`
}

type SParam struct {
	S11  Complex `json:"s11"`
	S12  Complex `json:"s12"`
	S21  Complex `json:"s21"`
	S22  Complex `json:"s22"`
	Freq uint64  `json:"freq"`
}

type Range struct {
	Start uint64 `json:"start"`
	End   uint64 `json:"end"`
}

type Command struct {
	ID      string `json:"id,omitEmpty"`
	Time    int    `json:"t,omitEmpty"`
	Command string `json:"cmd,omitEmpty"`
}

type RangeQuery struct {
	Command
	Range           Range        `json:"range"`
	Size            int          `json:"size"`
	LogDistribution bool         `json:"islog"`
	Avg             uint16       `json:"avg"`
	Select          SParamSelect `json:"sparam"`
	Result          []SParam     `json:"result,omitEmpty"`
}

// this command is not supported by pocket
// we have to handle this in the firmwarees
type CalibratedRangeQuery struct {
	Command
	Avg    uint16       `json:"avg"`
	Select SParamSelect `json:"sparam"`
	Result []SParam     `json:"result,omitEmpty"`
}

type SingleQuery struct {
	Command
	Freq   uint64       `json:"freq"`
	Avg    uint16       `json:"avg"`
	Select SParamSelect `json:"sparam"`
	Result SParam       `json:"result,omitEmpty"`
}

type ReasonableFrequencyRange struct {
	Command
	Result Range `json:"range"`
}

type Progress struct {
	Command
	Percentage int `json:"pc"`
}

type CustomResult struct {
	Message string `json:"message"`
	Command interface{}
}

type Complex struct {
	Real float64 `json:"real"`
	Imag float64 `json:"imag"`
}
