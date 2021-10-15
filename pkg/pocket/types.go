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
	S11 complex128
	S12 complex128
	S21 complex128
	S22 complex128
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
