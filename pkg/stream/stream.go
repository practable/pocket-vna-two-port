/*
Package stream connects to a websocket server and transfers JSON messages corresponding to the
types in pkg/pocket i.e. commands and results for pocketVNA

*/

package stream

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/timdrysdale/go-pocketvna/pkg/pocket"
	"github.com/timdrysdale/go-pocketvna/pkg/reconws"

	//"github.com/timdrysdale/go-pocketvna/pkg/calibration"
	//"github.com/timdrysdale/go-pocketvna/pkg/rfswitch"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func New(u string, ctx context.Context) Stream {

	request := make(chan interface{})
	response := make(chan interface{})

	r := reconws.New()

	go r.Reconnect(ctx, u)

	// We receive requests from user
	// i.e. reverse sense to our own services
	go PipeInterfaceToWs(request, r.In, ctx)
	go PipeWsToInterface(r.Out, response, ctx)
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

func RunDirect(u string, ctx context.Context) {

	r := reconws.New()

	v := pocket.NewVNA()

	command := make(chan interface{})

	result := make(chan interface{})

	go v.Run(command, result, ctx)

	go r.Reconnect(ctx, u)

	go PipeWsToInterface(r.In, command, ctx)

	go PipeInterfaceToWs(result, r.Out, ctx)

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
			}

			switch strings.ToLower(c.Command) {

			case "rq", "rangequery", "rc", "rangecal":

				s := pocket.RangeQuery{}

				err := json.Unmarshal([]byte(msg.Data), &s)

				if err != nil {
					log.WithField("error", err).Warning("Could not turn unmarshal JSON for RangeQuery (rq) command - invalid or missing parameters in JSON?")
				}

				out <- s

			case "crq", "calibratedrangequery":

				s := pocket.CalibratedRangeQuery{}

				err := json.Unmarshal([]byte(msg.Data), &s)

				if err != nil {
					log.WithField("error", err).Warning("Could not turn unmarshal JSON for CalibratedRangeQuery (rq) command - invalid or missing parameters in JSON?")
				}

				out <- s

			case "sq", "singlequery":

				s := pocket.SingleQuery{}

				err := json.Unmarshal([]byte(msg.Data), &s)

				if err != nil {
					log.WithField("error", err).Warning("Could not turn unmarshal JSON for SingleQuery (sq) command - invalid or missing parameters in JSON?")
				}

				out <- s

			case "rr", "reasonablefrequencyrange":

				s := pocket.ReasonableFrequencyRange{}

				err := json.Unmarshal([]byte(msg.Data), &s)

				if err != nil {
					log.WithField("error", err).Warning("Could not turn unmarshal JSON for ReasonableFrequencyRange (rr) command - invalid or missing parameters in JSON?")
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

func PipeWsToInterfaceCal(in chan reconws.WsMessage, out chan interface{}, ctx context.Context) {

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
			}

			switch strings.ToLower(c.Command) {

			case "rq", "rangequey":

				s := pocket.RangeQuery{}

				err := json.Unmarshal([]byte(msg.Data), &s)

				if err != nil {
					log.WithField("error", err).Warning("Could not turn unmarshal JSON for RangeQuery (rq) command - invalid or missing parameters in JSON?")
				}

				out <- s

			case "sq", "singlequery":

				s := pocket.SingleQuery{}

				err := json.Unmarshal([]byte(msg.Data), &s)

				if err != nil {
					log.WithField("error", err).Warning("Could not turn unmarshal JSON for SingleQuery (sq) command - invalid or missing parameters in JSON?")
				}

				out <- s

			case "rr", "reasonablefrequencyrange":

				s := pocket.ReasonableFrequencyRange{}

				err := json.Unmarshal([]byte(msg.Data), &s)

				if err != nil {
					log.WithField("error", err).Warning("Could not turn unmarshal JSON for ReasonableFrequencyRange (rr) command - invalid or missing parameters in JSON?")
				}

				out <- s
			}

		}

	}

}

func PipeWsToInterfaceSwitch(in chan reconws.WsMessage, out chan interface{}, ctx context.Context) {

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
			}

			switch strings.ToLower(c.Command) {

			case "rq", "rangequey":

				s := pocket.RangeQuery{}

				err := json.Unmarshal([]byte(msg.Data), &s)

				if err != nil {
					log.WithField("error", err).Warning("Could not turn unmarshal JSON for RangeQuery (rq) command - invalid or missing parameters in JSON?")
				}

				out <- s

			case "sq", "singlequery":

				s := pocket.SingleQuery{}

				err := json.Unmarshal([]byte(msg.Data), &s)

				if err != nil {
					log.WithField("error", err).Warning("Could not turn unmarshal JSON for SingleQuery (sq) command - invalid or missing parameters in JSON?")
				}

				out <- s

			case "rr", "reasonablefrequencyrange":

				s := pocket.ReasonableFrequencyRange{}

				err := json.Unmarshal([]byte(msg.Data), &s)

				if err != nil {
					log.WithField("error", err).Warning("Could not turn unmarshal JSON for ReasonableFrequencyRange (rr) command - invalid or missing parameters in JSON?")
				}

				out <- s
			}

		}

	}

}
