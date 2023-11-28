package calibration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/practable/pocket-vna-two-port/pkg/pocket"
	"github.com/practable/pocket-vna-two-port/pkg/reconws"
)

type Calibration struct {
	u        string
	R        *reconws.ReconWs
	Ctx      context.Context
	Request  chan interface{}
	Response chan interface{}
	Timeout  time.Duration
	Command  Command
	Scan     interface{}
	Ready    bool
}

/* Command object definition in python calibration service

   {
    "cmd":"twoportport",
    "freq": np.linspace(1e6,100e6,num=N),
    "short":{
        "s11": {
            "real":np.random.rand(N),
            "imag":np.random.rand(N),
            },
         "s12" : {
            "real":np.random.rand(N),
            "imag":np.random.rand(N),
            },
        "s21": {
            "real":np.random.rand(N),
            "imag":np.random.rand(N),
            },
         "s22" : {
            "real":np.random.rand(N),
            "imag":np.random.rand(N),
            },
      },
     "open":{ //as for short  },
     "load":{ // as for short },
     "load":{ // as for short },
     "dut":{ // as for short  }
    }
*/

// Command represents a twoport calibration command (request)
type Command struct {
	Command string   `json:"cmd"`
	Freq    []uint64 `json:"freq"`
	Short   Result   `json:"short"`
	Open    Result   `json:"open"`
	Load    Result   `json:"load"`
	Thru    Result   `json:"thru"`
	DUT     Result   `json:"dut"`
}

/* 2-port Result object definition in python calibration service
{
   "freq": network.f,
   "s11": {
      "real": np.squeeze(network.s_re[:,0,0]),
      "imag": np.squeeze(network.s_im[:,0,0]),
           },
  "s12" : {
      "real": np.squeeze(network.s_re[:,0,1]),
      "imag": np.squeeze(network.s_im[:,0,1]),
     },
 "s21": {
      "real": np.squeeze(network.s_re[:,1,0]),
      "imag": np.squeeze(network.s_im[:,1,0]),
     },
  "s22" : {
      "real": np.squeeze(network.s_re[:,1,1]),
      "imag": np.squeeze(network.s_im[:,1,1]),
     },
   }
*/

type Result struct {
	Freq []uint64     `json:"freq"`
	S11  ComplexArray `json:"s11"`
	S12  ComplexArray `json:"s12"`
	S21  ComplexArray `json:"s21"`
	S22  ComplexArray `json:"s22"`
}

type SParam struct {
	S11 ComplexArray `json:"s11"`
	S12 ComplexArray `json:"s12"`
	S21 ComplexArray `json:"s21"`
	S22 ComplexArray `json:"s22"`
}

type ComplexArray struct {
	Real []float64 `json:"real"`
	Imag []float64 `json:"imag"`
}

// MakeTwoPort assembles measurements into a twoport calibration command

func MakeTwoPort(pshort1, pshort2, popen1, popen2, pload1, pload2, pthru, pdut []pocket.SParam) (Command, error) {

	ls := []int{
		len(pshort1),
		len(pshort2),
		len(popen1),
		len(popen2),
		len(pload1),
		len(pload2),
		len(pthru),
		len(pdut),
	}

	for _, l := range ls {
		if l != ls[0] {
			return Command{}, errors.New("inputs are not the same length")
		}
	}

	short1 := PocketToResult(pshort1)
	short2 := PocketToResult(pshort2)
	open1 := PocketToResult(popen1)
	open2 := PocketToResult(popen2)
	load1 := PocketToResult(pload1)
	load2 := PocketToResult(pload2)
	thru := PocketToResult(pthru)
	dut := PocketToResult(pdut)

	short := CombineTwoReflectiveStandards(short1, short2)
	open := CombineTwoReflectiveStandards(open1, open2)
	load := CombineTwoReflectiveStandards(load1, load2)

	return Command{
		Command: "twoport",
		Freq:    short.Freq,
		Short:   short,
		Open:    open,
		Load:    load,
		Thru:    thru,
		DUT:     dut,
	}, nil

}

// CombineTwoReflectiveStandards combines S11 from the first argument,
// with S22 from the second, and sets S12, S21 to zero
func CombineTwoReflectiveStandards(r1, r2 Result) Result {
	z := make([]float64, len(r1.S11.Real))

	caz := ComplexArray{
		Real: z,
		Imag: z,
	}

	return Result{
		Freq: r1.Freq,
		S11:  r1.S11,
		S12:  caz,
		S21:  caz,
		S22:  r2.S22,
	}
}

// CombineTwoReflectiveStandards combines S11 from the first argument,
// with S22 from the second, and sets S12, S21 to zero
func CombineTwoReflectiveStandardsPocket(r1, r2 []pocket.SParam) []pocket.SParam {

	p := []pocket.SParam{}

	for i := range r1 {
		p = append(p, pocket.SParam{
			Freq: r1[i].Freq,
			S11:  r1[i].S11,
			S22:  r2[i].S22,
		})
	}

	return p

}

// PocketToCalibration changes the array structure to better suit our
// usages of the results (e.g. calibration calculations, and presentation
// in the user interface both need a frequency list so let's prepare that now
// for efficiency)
func PocketToCalibration(p []pocket.SParam) ([]uint64, SParam) {

	// we'll use append rather than assuming max-length array
	var freq []uint64
	var s11_real, s11_imag, s12_real, s12_imag, s21_real, s21_imag, s22_real, s22_imag []float64

	for _, param := range p {
		freq = append(freq, param.Freq)
		s11_real = append(s11_real, param.S11.Real)
		s11_imag = append(s11_imag, param.S11.Imag)
		s12_real = append(s12_real, param.S12.Real)
		s12_imag = append(s12_imag, param.S12.Imag)
		s21_real = append(s21_real, param.S21.Real)
		s21_imag = append(s21_imag, param.S21.Imag)
		s22_real = append(s22_real, param.S22.Real)
		s22_imag = append(s22_imag, param.S22.Imag)
	}

	ca := SParam{
		S11: ComplexArray{
			Real: s11_real,
			Imag: s11_imag,
		},
		S12: ComplexArray{
			Real: s12_real,
			Imag: s12_imag,
		},
		S21: ComplexArray{
			Real: s21_real,
			Imag: s21_imag,
		},
		S22: ComplexArray{
			Real: s22_real,
			Imag: s22_imag,
		},
	}
	return freq, ca

}

func PocketToResult(p []pocket.SParam) Result {
	freq, ca := PocketToCalibration(p)
	return Result{
		Freq: freq,
		S11:  ca.S11,
		S12:  ca.S12,
		S21:  ca.S21,
		S22:  ca.S22,
	}
}

// CalibrationToPocket is used by the middleware
func CalibrationToPocket(result Result) ([]pocket.SParam, error) {

	pa := []pocket.SParam{}

	if len(result.Freq) != len(result.S11.Real) { //assume real is same length as imag
		return pa, errors.New("Freq and S11 are different lengths")
	}
	if len(result.Freq) != len(result.S12.Real) {
		return pa, errors.New("Freq and S12 are different lengths")
	}
	if len(result.Freq) != len(result.S21.Real) {
		return pa, errors.New("Freq and S21 are different lengths")
	}
	if len(result.Freq) != len(result.S22.Real) {
		return pa, errors.New("Freq and S22 are different lengths")
	}

	for i, freq := range result.Freq {
		p := pocket.SParam{
			Freq: freq,
			S11: pocket.Complex{
				Real: result.S11.Real[i],
				Imag: result.S11.Imag[i],
			},
			S12: pocket.Complex{
				Real: result.S12.Real[i],
				Imag: result.S12.Imag[i],
			},
			S21: pocket.Complex{
				Real: result.S21.Real[i],
				Imag: result.S21.Imag[i],
			},
			S22: pocket.Complex{
				Real: result.S22.Real[i],
				Imag: result.S22.Imag[i],
			},
		}
		pa = append(pa, p)
	}

	return pa, nil
}

// New creates a Calibration object and the channels
// needed to communicate with the calibration server
// which are request, response
func New(u string, ctx context.Context) Calibration {

	request := make(chan interface{})
	response := make(chan interface{})

	r := reconws.New()

	go r.Reconnect(ctx, u)

	go PipeInterfaceToWs(request, r.Out, ctx)
	go PipeWsToInterface(r.In, response, ctx)

	c := Calibration{
		u:        u,
		R:        r,
		Ctx:      ctx,
		Request:  request,
		Response: response,
		Timeout:  time.Second,
		Command:  Command{},
		Scan:     pocket.RangeQuery{},
		Ready:    false,
	}

	c.Clear() //prepare for first use

	return c

}

// For this two-port experiment, we assume
// a twoport calibration every time
// Any corrections requiring the use of a one-port
// cal as part of that two port cal, are handled
// internal to the calibration service.
func (c *Calibration) Clear() {
	c.Command = Command{
		Command: "twoport",
	}
	c.Scan = pocket.RangeQuery{
		Select: pocket.SParamSelect{
			S11: true,
			S12: true,
			S21: true,
			S22: true,
		},
	}
	c.Ready = false

}

// Check ensures the results are consistent lengths
func (r *Result) Check() error {

	if len(r.S11.Real) != len(r.S11.Imag) {
		err := errors.New("S11 Real and Imag are different lengths")
		return err
	}
	if len(r.Freq) != len(r.S11.Real) {
		err := errors.New("Freq and S11 Real/Imag are different lengths")
		return err
	}

	if len(r.S11.Real) != len(r.S11.Imag) {
		err := errors.New("S11 Real and Imag are different lengths")
		return err
	}
	if len(r.Freq) != len(r.S11.Real) {
		err := errors.New("Freq and S11 Real/Imag are different lengths")
		return err
	}

	if len(r.S12.Real) != len(r.S12.Imag) {
		err := errors.New("S12 Real and Imag are different lengths")
		return err
	}
	if len(r.Freq) != len(r.S12.Real) {
		err := errors.New("Freq and S12 Real/Imag are different lengths")
		return err
	}

	if len(r.S21.Real) != len(r.S21.Imag) {
		err := errors.New("S11 Real and Imag are different lengths")
		return err
	}
	if len(r.Freq) != len(r.S11.Real) {
		err := errors.New("Freq and S11 Real/Imag are different lengths")
		return err
	}

	if len(r.S22.Real) != len(r.S22.Imag) {
		err := errors.New("S11 Real and Imag are different lengths")
		return err
	}
	if len(r.Freq) != len(r.S22.Real) {
		err := errors.New("Freq and S11 Real/Imag are different lengths")
		return err
	}

	return nil

}

// CompareFreq compares two lists using a threshold to judge equivalance.
// This is because there are some small differences in calculating
// intermediate frequencies in ranges on different hardware,
// we can run a check that accepts numerical differences that are too
// small to be scientifically relevant
func (r *Result) CompareFreq(freq []uint64) error {

	if len(r.Freq) != len(freq) {
		err := errors.New("Frequency arrays are different lengths")
		return err
	}

	thresh := uint64(100)

	for i, f := range r.Freq {

		if (f - freq[i]) >= thresh {
			err := fmt.Errorf("Frequency point %d differs by more than %d Hz (%d vs %d)", i, thresh, f, freq[i])
			return err
		}

	}

	return nil
}

// TODO check array format for full 2-port S-params
func (c *Calibration) SetShortParam(p []pocket.SParam) error {
	return c.SetShort(PocketToResult(p))
}

func (c *Calibration) SetOpenParam(p []pocket.SParam) error {
	return c.SetOpen(PocketToResult(p))
}

func (c *Calibration) SetLoadParam(p []pocket.SParam) error {
	return c.SetLoad(PocketToResult(p))
}

func (c *Calibration) SetThruParam(p []pocket.SParam) error {
	return c.SetThru(PocketToResult(p))
}

func (c *Calibration) SetDUTParam(p []pocket.SParam) error {
	return c.SetDUT(PocketToResult(p))
}

// SetShort adds the measurement for the standard short to the cal object
func (c *Calibration) SetShort(result Result) error {
	return c.Set("short", result)
}

//SetOpen adds the measurement for the standard open to the cal object
func (c *Calibration) SetOpen(result Result) error {
	return c.Set("open", result)
}

//SetLoad adds the measurement for the standard load to the cal object
func (c *Calibration) SetLoad(result Result) error {
	return c.Set("load", result)
}

//SetThru adds the measurement for the standard thru to the cal object
func (c *Calibration) SetThru(result Result) error {
	return c.Set("thru", result)
}

// SetDUT adds the measurement for the DUT to the cal object
func (c *Calibration) SetDUT(result Result) error {
	return c.Set("dut", result)
}

// Set adds results to the cal object, ensuring that the frequency
// field is populated
func (c *Calibration) Set(target string, result Result) error {

	err := result.CompareFreq(c.Command.Freq)

	if err != nil {
		// if freq is empty, then assume cleared command
		if len(c.Command.Freq) == 0 {
			c.Command.Freq = result.Freq
		} else {
			return err
		}
	}

	switch {
	case target == "short":
		c.Command.Short = result
	case target == "open":
		c.Command.Open = result
	case target == "load":
		c.Command.Load = result
	case target == "thru":
		c.Command.Thru = result
	case target == "dut":
		c.Command.DUT = result
	default:
		return fmt.Errorf("Unknown measurement: %s. Should be short, open, load, thru, or dut.", target)
	}

	return nil

}

// no need for errors with messages - should have found these out by now -
// this is mainly a check that all required elements are present
// but adding the belt and braces of checking individual array lengths
func (ca *ComplexArray) BadLen(freq []uint64) bool {

	if len(ca.Real) != len(ca.Imag) {
		return true
	}

	if len(freq) != len(ca.Real) {
		return true
	}

	return false
}

// no need for errors with messages - should have found these out by now -
// this is mainly a check that all required elements are present
// but adding the belt and braces of checking individual array lengths
func (r *Result) LengthEquals(required int) bool {

	rl := []int{
		len(r.Freq),
		len(r.S11.Real),
		len(r.S11.Imag),
		len(r.S12.Real),
		len(r.S12.Imag),
		len(r.S21.Real),
		len(r.S21.Imag),
		len(r.S22.Real),
		len(r.S22.Imag),
	}

	for _, l := range rl {
		if l != required {
			return false
		}
	}

	return true
}

func (c *Calibration) Apply() (Result, error) {

	result := Result{}
	required := len(c.Command.Freq)

	if !c.Command.Short.LengthEquals(required) {
		return result, errors.New("Wrong number of samples in Short data")
	}

	if !c.Command.Open.LengthEquals(required) {
		return result, errors.New("Wrong number of samples in Open data")
	}

	if !c.Command.Load.LengthEquals(required) {
		return result, errors.New("Wrong number of samples in Load data")
	}
	if !c.Command.Thru.LengthEquals(required) {
		return result, errors.New("Wrong number of samples in Thru data")
	}
	if !c.Command.DUT.LengthEquals(required) {
		return result, errors.New("Wrong number of samples in DUT data")
	}

	select {
	case c.Request <- c.Command:
		// carry on
	case <-time.After(c.Timeout):
		return result, errors.New("Timeout on sending request to calibration service")
	}

	select {
	case response := <-c.Response:

		result, ok := response.(Result)

		if ok {
			return result, nil
		} else {
			log.Warn(fmt.Sprintf("%+v", response))
			return result, errors.New("Response did not contain a result")
		}

	case <-time.After(c.Timeout):
		return result, errors.New("Timeout receiving response from calibration service")
	}

}

//works for any serialisable struct
func PipeInterfaceToWs(in chan interface{}, out chan reconws.WsMessage, ctx context.Context) {

	mtype := int(websocket.TextMessage)

	for {
		select {

		case <-ctx.Done():
			return
		case s := <-in:

			payload, err := json.Marshal(s)

			if err != nil {
				log.WithField("error", err).Warning("Could not turn interface{} into JSON")
			}

			out <- reconws.WsMessage{Data: payload, Type: mtype}

		}

	}

}

// only works for defined structs
func PipeWsToInterface(in chan reconws.WsMessage, out chan interface{}, ctx context.Context) {

	for {
		select {

		case <-ctx.Done():
			return

		case msg := <-in:

			var r Result

			err := json.Unmarshal([]byte(msg.Data), &r)

			if err != nil {
				log.WithField("error", err).Warning("Could not turn unmarshal JSON - invalid report string in JSON?")
			}

			out <- r

		}

	}

}
