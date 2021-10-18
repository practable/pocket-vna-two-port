/*
Package stream connects to a websocket server and transfers JSON messages corresponding to the
types in pkg/pocket i.e. commands and results for pocketVNA

*/

package stream

import (
	"context"
	"encoding/json"
	"strings"

	"../pocket"
	"../reconws"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func Run(u string, ctx context.Context) {

	r := reconws.New()

	v := pocket.NewVNA()

	command := make(chan interface{})

	result := make(chan interface{})

	go v.Run(command, result, ctx)

	go r.Reconnect(ctx, u)

	go PipeWsToInterface(r.In, command, ctx)

	go PipeInterfaceToWs(result, r.Out, ctx)

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
