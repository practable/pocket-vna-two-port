package calibration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/timdrysdale/go-pocketvna/pkg/pocket"
	"github.com/timdrysdale/go-pocketvna/pkg/reconws"
)

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

func (c *Calibration) Apply() (Result, error) {

	result := Result{}

	if c.Command.Short.BadLen(c.Command.Freq) {
		return result, errors.New("Wrong number of samples in Short data")
	}

	if c.Command.Open.BadLen(c.Command.Freq) {
		return result, errors.New("Wrong number of samples in Open data")
	}

	if c.Command.Load.BadLen(c.Command.Freq) {
		return result, errors.New("Wrong number of samples in Load data")
	}
	if c.Command.Thru.BadLen(c.Command.Freq) {
		return result, errors.New("Wrong number of samples in Thru data")
	}
	if c.Command.DUT.BadLen(c.Command.Freq) {
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
