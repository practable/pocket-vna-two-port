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
	Freq uint64
	S11  complex128
	S12  complex128
	S21  complex128
	S22  complex128
}

type Name int

const (
	RangeQuery = iota
	SingleQuery
	GetReasonableFrequencyRange
)

type Command struct {
	Name            Name
	Freq            uint64
	Start           uint64
	End             uint64
	Size            int
	LogDistribution bool
	Avg             uint16
	Select          SParamSelect
}

type Data struct {
	SParams []Sparam
	Select  SParamSelect
	Start   uint64
	End     uint64
}
