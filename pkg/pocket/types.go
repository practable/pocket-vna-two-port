package pocket

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
	Serial       string
	Manufacturer string
	Product      string
	Release      int
}

type SParamSelect struct {
	S11 bool
	S12 bool
	S21 bool
	S22 bool
}

type SParam struct {
	S11  Complex
	S12  Complex
	S21  Complex
	S22  Complex
	Freq uint64
}

type Range struct {
	Start uint64
	End   uint64
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
	LogDistribution bool         `json:"isLog"`
	Avg             uint16       `json:"avg"`
	Select          SParamSelect `json:"sparam"`
	Result          []SParam     `json:"result,omitEmpty"`
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
	Message string
	Command interface{}
}

type Complex struct {
	Real float64
	Imag float64
}
