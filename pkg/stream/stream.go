/*
Package stream connects to a websocket server and transfers JSON messages corresponding to the
types in pkg/pocket i.e. commands and results for pocketVNA
This is done here to simplify the message handling in middleware package `middle`

*/

package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/practable/pocket-vna-two-port/pkg/pocket"
	"github.com/practable/pocket-vna-two-port/pkg/reconws"
	log "github.com/sirupsen/logrus"
)

type Stream struct {
	u        string
	R        *reconws.ReconWs
	Ctx      context.Context
	Request  chan interface{}
	Response chan interface{}
	Timeout  time.Duration
}

// TODO duplicate the testing applied to RunDirect
func New(ctx context.Context, u string) Stream {

	request := make(chan interface{}, 2)
	response := make(chan interface{}, 2)

	r := reconws.New()

	go r.Reconnect(ctx, u)

	// We receive requests from user
	// i.e. reverse sense to our own services

	go PipeWsToInterface(r.In, request, ctx)

	go PipeInterfaceToWs(response, r.Out, ctx)

	go HeartBeat(r.Out, time.Second, ctx)

	return Stream{
		u:        u,
		R:        r,
		Ctx:      ctx,
		Request:  request,
		Response: response,
		Timeout:  time.Second,
	}

}

// This is the straight-forward version of the firmware with no added functionality
// useful for raw access to the VNA. Requires VNA to be connected for testing.
func RunDirect(ctx context.Context, u string) {

	r := reconws.New()

	h, disconnect, err := pocket.NewHardware()

	if err != nil {
		log.Errorf("Stream.RunDirect: %s", err.Error())
		return
	}

	go func() {
		<-ctx.Done()
		disconnect()
	}()

	go func() {

		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-r.In:

				var c pocket.Command

				err := json.Unmarshal([]byte(msg.Data), &c)

				if err != nil {
					log.WithField("err", err).Warning("Could not turn unmarshal JSON - invalid cmd string in JSON?")
					fmt.Printf("\n%s\n", msg.Data)
				}

				log.Debugf("pkg/stream.RunDirect() received command %s", strings.ToLower(c.Command))

				switch strings.ToLower(c.Command) {

				case "rq", "rangequery", "rc", "rangecal", "sc", "setupcal", "mc", "measurecal", "cc", "confirmcal":

					s := pocket.RangeQuery{}

					err := json.Unmarshal([]byte(msg.Data), &s)

					if err != nil {
						log.WithField("err", err).Warning("Could not turn unmarshal JSON for RangeQuery (rq) command - invalid or missing parameters in JSON?")
						fmt.Printf("\n%s\n", msg.Data)
					}

					err = h.HandleCommand(s)

					if err != nil {

						log.WithFields(log.Fields{"err": err.Error(), "command": s}).Error("error handling command")

						cr := pocket.CustomResult{Message: err.Error(), Command: s}

						data, err := json.Marshal(cr)

						if err != nil {
							data = []byte("could not marshal error message")
						}

						r.Out <- reconws.WsMessage{Data: data, Type: msg.Type}

						continue

					}

					data, err := json.Marshal(s)

					if err != nil {

						log.WithFields(log.Fields{"err": err.Error(), "command": s}).Error("error marshalling result")

						cr := pocket.CustomResult{Message: err.Error(), Command: s}

						data, err = json.Marshal(cr)

						if err != nil {
							data = []byte("could not marshal error message")
						}

					}

					r.Out <- reconws.WsMessage{Data: data, Type: msg.Type}

				case "crq", "calibratedrangequery":

					s := pocket.CalibratedRangeQuery{}

					err := json.Unmarshal([]byte(msg.Data), &s)

					if err != nil {
						log.WithField("err", err).Warning("Could not turn unmarshal JSON for CalibratedRangeQuery (rq) command - invalid or missing parameters in JSON?")
						fmt.Printf("\n%s\n", msg.Data)
					}

				case "sq", "singlequery":

					s := pocket.SingleQuery{}

					err := json.Unmarshal([]byte(msg.Data), &s)

					if err != nil {
						log.WithField("err", err).Warning("Could not turn unmarshal JSON for SingleQuery (sq) command - invalid or missing parameters in JSON?")
						fmt.Printf("\n%s\n", msg.Data)
					}

					err = h.HandleCommand(s)

					if err != nil {

						log.WithFields(log.Fields{"err": err.Error(), "command": s}).Error("error handling command")

						cr := pocket.CustomResult{Message: err.Error(), Command: s}

						data, err := json.Marshal(cr)

						if err != nil {
							data = []byte("could not marshal error message")
						}

						r.Out <- reconws.WsMessage{Data: data, Type: msg.Type}

						continue

					}

					data, err := json.Marshal(s)

					if err != nil {

						log.WithFields(log.Fields{"err": err.Error(), "command": s}).Error("error marshalling result")

						cr := pocket.CustomResult{Message: err.Error(), Command: s}

						data, err = json.Marshal(cr)

						if err != nil {
							data = []byte("could not marshal error message")
						}

					}

					r.Out <- reconws.WsMessage{Data: data, Type: msg.Type}

				case "rr", "reasonablefrequencyrange":

					s := pocket.ReasonableFrequencyRange{}

					err := json.Unmarshal([]byte(msg.Data), &s)

					if err != nil {
						log.WithField("err", err).Warning("Could not turn unmarshal JSON for ReasonableFrequencyRange (rr) command - invalid or missing parameters in JSON?")
						fmt.Printf("\n%s\n", msg.Data)
					}

					err = h.HandleCommand(s)

					if err != nil {

						log.WithFields(log.Fields{"err": err.Error(), "command": s}).Error("error handling command")

						cr := pocket.CustomResult{Message: err.Error(), Command: s}

						data, err := json.Marshal(cr)

						if err != nil {
							data = []byte("could not marshal error message")
						}

						r.Out <- reconws.WsMessage{Data: data, Type: msg.Type}

						continue

					}

					data, err := json.Marshal(s)

					if err != nil {

						log.WithFields(log.Fields{"err": err.Error(), "command": s}).Error("error marshalling result")

						cr := pocket.CustomResult{Message: err.Error(), Command: s}

						data, err = json.Marshal(cr)

						if err != nil {
							data = []byte("could not marshal error message")
						}

					}

					r.Out <- reconws.WsMessage{Data: data, Type: msg.Type}

				}

			}
		}

	}()

	go r.Reconnect(ctx, u)

	go HeartBeat(r.Out, time.Second, ctx)

}

func HeartBeat(out chan reconws.WsMessage, t time.Duration, ctx context.Context) {

	mtype := int(websocket.TextMessage)

	for {
		select {

		case <-ctx.Done():
			return
		case <-time.After(t):
			out <- reconws.WsMessage{
				Data: []byte("{\"cmd\":\"hb\"}"),
				Type: mtype,
			}

		}
	}

}

func PipeWsToInterface(in chan reconws.WsMessage, out chan interface{}, ctx context.Context) {

	for {
		select {

		case <-ctx.Done():
			return

		case msg := <-in:

			//var rq pocket.RangeQuery
			//var sq pocket.SingleQuery
			//var rr pocket.ReasonableFrequencyRange

			var c pocket.Command

			err := json.Unmarshal([]byte(msg.Data), &c)

			if err != nil {
				log.WithField("error", err).Warning("Could not turn unmarshal JSON - invalid cmd string in JSON?")
				fmt.Printf("\n%s\n", msg.Data)
			}

			switch strings.ToLower(c.Command) {

			case "rq", "rangequery", "rc", "rangecal", "sc", "setupcal", "mc", "mesaurecal", "cc", "confirmcal":

				s := pocket.RangeQuery{}

				err := json.Unmarshal([]byte(msg.Data), &s)

				if err != nil {
					log.WithField("error", err).Warning("Could not turn unmarshal JSON for RangeQuery (rq) command - invalid or missing parameters in JSON?")
					fmt.Printf("\n%s\n", msg.Data)
				}

				out <- s

			case "crq", "calibratedrangequery":

				s := pocket.CalibratedRangeQuery{}

				err := json.Unmarshal([]byte(msg.Data), &s)

				if err != nil {
					log.WithField("error", err).Warning("Could not turn unmarshal JSON for CalibratedRangeQuery (rq) command - invalid or missing parameters in JSON?")
					fmt.Printf("\n%s\n", msg.Data)
				}

				out <- s

			case "sq", "singlequery":

				s := pocket.SingleQuery{}

				err := json.Unmarshal([]byte(msg.Data), &s)

				if err != nil {
					log.WithField("error", err).Warning("Could not turn unmarshal JSON for SingleQuery (sq) command - invalid or missing parameters in JSON?")
					fmt.Printf("\n%s\n", msg.Data)
				}

				out <- s

			case "rr", "reasonablefrequencyrange":

				s := pocket.ReasonableFrequencyRange{}

				err := json.Unmarshal([]byte(msg.Data), &s)

				if err != nil {
					log.WithField("error", err).Warning("Could not turn unmarshal JSON for ReasonableFrequencyRange (rr) command - invalid or missing parameters in JSON?")
					fmt.Printf("\n%s\n", msg.Data)
				}

				out <- s
			}

		}

	}

}

// This can be used for all of the external connections because it is data structure agnostic
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
