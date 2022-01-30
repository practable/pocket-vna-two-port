package rfswitch

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

func New(u string, ctx context.Context) Switch {

	request := make(chan interface{})
	response := make(chan interface{})

	r := reconws.New()

	go r.Reconnect(ctx, u)

	go PipeInterfaceToWs(request, r.Out, ctx)
	go PipeWsToInterface(r.In, response, ctx)

	return Switch{
		u:        u,
		R:        r,
		Ctx:      ctx,
		Request:  request,
		Response: response,
		Timeout:  time.Second,
	}
}

func (s *Switch) SetShort() error {
	return s.SetPort("short")
}

func (s *Switch) SetOpen() error {
	return s.SetPort("open")
}

func (s *Switch) SetLoad() error {
	return s.SetPort("load")
}

func (s *Switch) SetDUT() error {
	return s.SetPort("dut")
}

func (s *Switch) SetPort(port string) error {
	request := Command{
		Set: "port",
		To:  port,
	}

	select {
	case <-time.After(s.Timeout):
		return errors.New("timeout sending request")
	case s.Request <- request:
		//carry on
	}

	select {
	case <-time.After(s.Timeout):
		return errors.New("timeout receiving response")
	case response := <-s.Response:
		r, ok := response.(Report)

		if !ok {
			return errors.New("Unexpected response")
		}

		if r.Report == "error" {
			return errors.New("Error" + r.Is)
		}

		if r.Report == "port" {

			if r.Is == port {
				return nil
			} else {
				return errors.New("Wrong port set")
			}
		}

		// catch anything else
		return errors.New("unexpected response")

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

func PipeWsToInterface(in chan reconws.WsMessage, out chan interface{}, ctx context.Context) {

	for {
		select {

		case <-ctx.Done():
			return

		case msg := <-in:

			var r Report

			err := json.Unmarshal(msg.Data, &r)

			if err != nil {
				log.WithField("error", err).Warning(fmt.Sprintf("Could not turn unmarshal JSON - invalid report string in JSON? %s", string(msg.Data)))
			}

			out <- r

		}

	}

}
