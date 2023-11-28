package rfswitch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/practable/pocket-vna-two-port/pkg/reconws"
)

func New(u string, ctx context.Context) Switch {

	request := make(chan interface{})
	response := make(chan interface{})

	r := reconws.New()

	go r.Reconnect(ctx, u)

	go PipeInterfaceToWs(request, r.Out, ctx)
	go PipeWsToInterface(r.In, response, ctx)

	return Switch{
		u:            u,
		R:            r,
		Ctx:          ctx,
		Request:      request,
		Response:     response,
		Timeout:      2 * time.Second,
		DrainTimeout: 10 * time.Millisecond,
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

func (s *Switch) SetThru() error {
	return s.SetPort("thru")
}
func (s *Switch) SetDUT1() error {
	return s.SetPort("dut1")
}
func (s *Switch) SetDUT2() error {
	return s.SetPort("dut2")
}
func (s *Switch) SetDUT3() error {
	return s.SetPort("dut3")
}
func (s *Switch) SetDUT4() error {
	return s.SetPort("dut4")
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
		log.Infof("rfswitch: sent command to pipe %s", port)
		//carry on
	}

	for i := 0; i < 5; i++ {

		select {
		case <-time.After(s.Timeout):
			log.Infof("timeout setting rfswitch to %s on %d loop", port, i)
			return errors.New("timeout receiving response")

		case response := <-s.Response:

			log.Infof("rfswitch: got a response, in raw is %+v", s.Response)

			r, ok := response.(Report)

			if ok {

				if r.Report == "error" {
					log.Infof("rfswitch returned error when setting to port %s of %s", port, r.Is)
					return errors.New("Error" + r.Is)
				}

				if r.Report == "port" && r.Is == port {
					log.Info("rfswitch: got the response we wanted")
					return nil
				}
				log.Infof("rfswitch set to wrong port, wanted %s got %s on %d loop", r.Is, port, i)
				// if get to here, then we have a valid response
				// but with the wrong port, and we'll ignore it
				// else we throw errors forever after getting one timeout.
				// Just wait to see if a valid response is given in the
				// right time frame.
				// To avoid false positives, we could number requests and responses.

			}
			log.Infof("rfswitch - ignoring this reply, probably a blank line")
			// not a report message - probably a blank line, ignore
		}
	}

	// if we get to here, too many blank lines or non-standard
	// messages were sent - check arduino software and USB connection?
	return errors.New("Too many Unexpected responses")

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
			log.Infof("rfswitch: hardware replied: %s", string(msg.Data))
			out <- r
			log.Tracef("rfswitch: wrote hardware reply to pipe: %s", string(msg.Data))
		}

	}

}
