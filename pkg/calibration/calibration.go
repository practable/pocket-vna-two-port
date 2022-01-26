package calibration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/timdrysdale/go-pocketvna/pkg/reconws"
)

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

// we only have one calibration method for now
// so we can set the command up front
func (c *Calibration) Clear() {
	c.Command = Command{
		Command: "oneport",
	}
}

func (r *Result) Check() error {

	if len(r.S11.Real) != len(r.S11.Imag) {
		err := errors.New("Real and Imag are different lengths")
		return err
	}

	if len(r.Freq) != len(r.S11.Real) {
		err := errors.New("Freq and Real/Imag are different lengths")
		return err
	}

	return nil

}

// Just in case there are some odd calculation errors in generating
// the frequency points between runs

func UnequalFreq(a, b, thresh uint64) bool {

	return ((a - b) >= thresh)

}

func (r *Result) CompareFreq(freq []uint64) error {

	if len(r.Freq) != len(freq) {
		err := errors.New("Frequency arrays are different lengths")
		return err
	}

	thresh := uint64(100)

	for i, f := range r.Freq {

		if UnequalFreq(f, freq[i], thresh) {
			err := fmt.Errorf("Frequency point %d differs by more than %d Hz (%d vs %d)", i, thresh, f, freq[i])
			return err
		}

	}

	return nil
}

func (c *Calibration) SetShort(result Result) error {
	return c.Set("short", result)
}

func (c *Calibration) SetOpen(result Result) error {
	return c.Set("open", result)
}
func (c *Calibration) SetLoad(result Result) error {
	return c.Set("load", result)
}

func (c *Calibration) SetDUT(result Result) error {
	return c.Set("dut", result)
}

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
		c.Command.Short = result.S11
	case target == "open":
		c.Command.Open = result.S11
	case target == "load":
		c.Command.Load = result.S11
	case target == "dut":
		c.Command.DUT = result.S11
	default:
		return fmt.Errorf("Unknown measurement: %s. Should be short, open, load, or dut.", target)
	}

	return nil

}

// no need for errors with messages - should have found these out by now -
// this is mainly a check that all requirement elements are present
// but adding the belt and braces of checking individual array lengths
func (ca *ComplexArray) BadLen(freq []uint64) bool {

	if len(ca.Real) != len(ca.Imag) {
		return false
	}

	if len(freq) != len(ca.Real) {
		return false
	}

	return true
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
