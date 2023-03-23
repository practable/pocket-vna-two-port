package pocket

import (
	"errors"
	"math"
	"reflect"

	log "github.com/sirupsen/logrus"
)

type VNA interface {
	Connect() (func() error, error)
	GetReasonableFrequencyRange(command interface{}) error
	HandleCommand(command interface{}) error
	RangeQuery(command interface{}) error
	SingleQuery(command interface{}) error
}

//TODO mock that takes a pointer to mock switch state, and returns different results depending on switch state

type Mock struct {
	ConnectError                   error
	DisconnectError                error
	CommandError                   error
	ResultRangeQuery               []SParam
	ResultSingleQuery              SParam
	ResultReasonableFrequencyRange Range
	CommandsReceived               []interface{}
}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) Connect() (func() error, error) {
	return func() error { return m.DisconnectError }, m.ConnectError
}

func (m *Mock) GetReasonableFrequencyRange(command interface{}) error {

	c := command.(*ReasonableFrequencyRange)

	c.Result.Start = m.ResultReasonableFrequencyRange.Start
	c.Result.End = m.ResultReasonableFrequencyRange.End

	cc := *c

	m.CommandsReceived = append(m.CommandsReceived, cc)

	command = c

	return m.CommandError
}

func (m *Mock) SingleQuery(command interface{}) error {

	c := command.(*SingleQuery)

	c.Result = m.ResultSingleQuery
	cc := *c
	m.CommandsReceived = append(m.CommandsReceived, cc)

	command = c

	return m.CommandError
}

func (m *Mock) RangeQuery(command interface{}) error {

	c := command.(*RangeQuery)

	c.Result = m.ResultRangeQuery
	cc := *c
	m.CommandsReceived = append(m.CommandsReceived, cc)

	command = c

	return m.CommandError
}

func (m *Mock) HandleCommand(command interface{}) error {

	// used to return CustomResult{Message: err.Error()} on error, or copy of command
	log.WithField("type", reflect.TypeOf(command)).Debugf("Handling Command")
	switch (command).(type) {

	case *ReasonableFrequencyRange:

		return m.GetReasonableFrequencyRange(command)

	case *RangeQuery:

		return m.RangeQuery(command)

	case *SingleQuery:

		return m.SingleQuery(command) //.(SingleQuery))

	default:
		return errors.New("unknown command")
	}

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
	What            string       `json:"what"`
}

// GenericCommand (for getting the command type in HandleCommand)
type GenericCommand struct {
	Command
}

// this command is not supported by pocket
// we have to handle this in the firmwarees
type CalibratedRangeQuery struct {
	Command
	What   string       `json:"what"`
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

const (
	Undefined Distribution = iota //handle default value being undefined
	Linear
	Log
)

type Distribution int

func NewHardware() (VNA, func() error, error) {
	h := new(Hardware)
	disconnect, err := h.Connect()

	if err != nil {
		return h, disconnect, err
	}

	return h, disconnect, nil
}

/* Run provides a go channel interface to the first available instance of a pocket VNA device

There are two uni-directional channels, one to receive commands, the other to reply with data

*/

// Connect returns a disconnect function that should called when finished
func (h *Hardware) Connect() (func() error, error) {

	handle, err := getFirstDeviceHandle()

	if err != nil {
		return func() error { return nil }, err
	}

	h.handle = handle

	return func() error { return releaseHandle(h.handle) }, nil
}

func ForceUnlockDevices() error {

	return forceUnlockDevices()

}

func (h *Hardware) GetReasonableFrequencyRange(command interface{}) error {

	r := command.(*ReasonableFrequencyRange)

	fStart, fEnd, err := getReasonableFrequencyRange(h.handle)

	if err != nil {
		return err
	}

	r.Result.Start = fStart
	r.Result.End = fEnd

	command = r

	return err

}

func (h *Hardware) HandleCommand(command interface{}) error {

	// used to return CustomResult{Message: err.Error()} on error, or copy of command
	switch (command).(type) {

	case *ReasonableFrequencyRange:

		return h.GetReasonableFrequencyRange(command)

	case *RangeQuery:

		return h.RangeQuery(command)

	case *SingleQuery:

		return h.SingleQuery(command) //.(SingleQuery))

	default:
		return errors.New("unknown command")
	}

}

func (h *Hardware) RangeQuery(command interface{}) error {

	r := command.(*RangeQuery)

	distr := 1 // Linear

	if r.LogDistribution {
		distr = 2
	}

	sparams, err := rangeQuery(h.handle, r.Range.Start, r.Range.End, r.Size, distr, r.Avg, r.Select)

	if err != nil {
		return err
	}

	r.Result = sparams

	command = r

	return err
}

func (h *Hardware) SingleQuery(command interface{}) error {

	s := command.(*SingleQuery)

	sparam, err := singleQuery(h.handle, s.Freq, s.Avg, s.Select)

	if err != nil {
		return err
	}

	s.Result = sparam

	command = s

	return err

}

func LinFrequency(start, end uint64, size int) []uint64 {

	var ff []uint64
	s := float64(start)
	e := float64(end)

	for i := 0; i < size; i++ {
		f := s + float64(i)*(e-s)/(float64(size)-1)
		ff = append(ff, uint64(f))
	}
	return ff
}

func LogFrequency(start, end uint64, size int) []uint64 {

	var ff []uint64
	s := float64(start)
	e := float64(end)
	x := e / s
	for i := 0; i < size; i++ {

		y := float64(i) / (float64(size) - 1.0)
		f := s * math.Pow(x, y)
		ff = append(ff, uint64(math.Round(f)))
	}
	return ff

}
