// Package middle coordinates the response to user requests that require the use of the rfswitch and calibration services
package middle

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/timdrysdale/go-pocketvna/pkg/calibration"
	"github.com/timdrysdale/go-pocketvna/pkg/pocket"
	"github.com/timdrysdale/go-pocketvna/pkg/rfswitch"
	"github.com/timdrysdale/go-pocketvna/pkg/stream"
)

func New(uc, ur, us string, ctx context.Context) Middle {

	c := calibration.New(uc, ctx)
	r := rfswitch.New(ur, ctx)
	s := stream.New(us, ctx)

	v := pocket.New(ctx)

	timeout := time.Second

	requesttimeout := 2 * time.Minute

	go func() {
		for {

			select {

			case request := <-s.Request:

				select {

				case s.Response <- HandleRequest(request, c, r, v):
					// carry on
				case <-time.After(requesttimeout):
					s.Response <- TimeoutMessage(request)
				}

			case <-time.After(timeout):
				//carry on
			case <-ctx.Done():
				return
			}

		} //for
	}() //func

	return Middle{
		Calibration: c,
		Stream:      s,
		Switch:      r,
		VNA:         v,
	}

}

func TimeoutMessage(request interface{}) interface{} {

	return pocket.CustomResult{
		Message: "timeout waiting for request to be handled",
		Command: request,
	}

}

func HandleRequest(request interface{}, c calibration.Calibration, r rfswitch.Switch, v pocket.VNAService) interface{} {

	switch request.(type) {

	case pocket.ReasonableFrequencyRange, pocket.SingleQuery:

		v.Request <- request

		return <-v.Response

	case pocket.RangeQuery:

		// this type is used for different commands

		rq := request.(pocket.RangeQuery)

		switch rq.Command.Command {

		case "rq", "rangequery":

			v.Request <- request
			return <-v.Response

		case "rc", "rangecal":

			log.WithFields(log.Fields{
				"request": rq,
			}).Infof("Middle.HandleRequest with ID: %s", rq.ID)

			return RangeCal(rq, c, r, v)

		default:
			return pocket.CustomResult{
				Message: "Unknown request",
				Command: request,
			}
		}

	case pocket.CalibratedRangeQuery:

		crq := request.(pocket.CalibratedRangeQuery)

		return CalibratedRangeQuery(crq, c, r, v)

	default:
		return pocket.CustomResult{
			Message: "Unknown request",
			Command: request,
		}
	}
}

func CalibratedRangeQuery(rq pocket.CalibratedRangeQuery, c calibration.Calibration, r rfswitch.Switch, v pocket.VNAService) interface{} {

	//TODO implement the application of the calibration

	scan := rq

	scan.Command.Command = "rq"

	v.Request <- scan

	return <-v.Response

}

func RangeCal(rc pocket.RangeQuery, c calibration.Calibration, r rfswitch.Switch, v pocket.VNAService) interface{} {

	// Check port 1 is specified

	onlyS11 := rc.Select.S11 && !rc.Select.S12 && !rc.Select.S21 && !rc.Select.S22

	if !onlyS11 {
		msg := fmt.Sprintf("Error: calibration is only supported on Port1 (S11). Resend the command with only S11 selected (true). You had S11:%v, S12:%v, S21:%v, S22:%v",
			rc.Select.S11, rc.Select.S12, rc.Select.S21, rc.Select.S22)
		return pocket.CustomResult{
			Message: msg,
			Command: rc,
		}
	}

	// clear previous cal
	c.Clear()

	// prepare the scanning command used to measure each standard
	scan := rc
	scan.Command.Command = "rq"

	//save it for the cqr to use later
	c.Scan = scan

	// SHORT

	name := "short"

	err := r.SetShort()

	if err != nil {
		return pocket.CustomResult{
			Message: "Error setting RF switch to " + name + ": " + err.Error(),
			Command: rc,
		}
	}

	v.Request <- scan

	response := <-v.Response

	rrq, ok := response.(pocket.RangeQuery)

	if !ok {
		return pocket.CustomResult{
			Message: "Error measuring " + name,
			Command: response,
		}
	}

	result := rrq.Result

	if len(result) != rc.Size {
		return pocket.CustomResult{
			Message: "Error measuring " + name,
			Command: response,
		}
	}

	err = c.SetShortParam(result)

	if err != nil {
		return pocket.CustomResult{
			Message: "Error putting data for " + name + " into cal store: " + err.Error(),
			Command: result,
		}
	}

	// OPEN

	name = "open"

	err = r.SetOpen()

	if err != nil {
		return pocket.CustomResult{
			Message: "Error setting RF switch to " + name + ": " + err.Error(),
			Command: rc,
		}
	}

	v.Request <- scan

	response = <-v.Response

	rrq, ok = response.(pocket.RangeQuery)

	if !ok {
		return pocket.CustomResult{
			Message: "Error measuring " + name,
			Command: response,
		}
	}

	result = rrq.Result

	if len(result) != rc.Size {
		return pocket.CustomResult{
			Message: "Error measuring " + name,
			Command: response,
		}
	}

	err = c.SetOpenParam(result)

	if err != nil {
		return pocket.CustomResult{
			Message: "Error putting data for " + name + " into cal store: " + err.Error(),
			Command: result,
		}
	}

	// LOAD

	name = "load"

	err = r.SetLoad()

	if err != nil {
		return pocket.CustomResult{
			Message: "Error setting RF switch to " + name + ": " + err.Error(),
			Command: rc,
		}
	}

	v.Request <- scan

	response = <-v.Response

	rrq, ok = response.(pocket.RangeQuery)

	if !ok {
		return pocket.CustomResult{
			Message: "Error measuring " + name,
			Command: response,
		}
	}

	result = rrq.Result

	if len(result) != rc.Size {
		return pocket.CustomResult{
			Message: "Error measuring " + name,
			Command: response,
		}
	}

	err = c.SetLoadParam(result)

	if err != nil {
		return pocket.CustomResult{
			Message: "Error putting data for " + name + " into cal store: " + err.Error(),
			Command: result,
		}
	}

	// send some results back so the success is confirmed with the presence of data
	rc.Result = rrq.Result

	return rc

	// don't return a custom result because with a command, with results, because
	// the json parser can't cope with this, causing failed tests. i.e. AVOID THIS:
	// return pocket.CustomResult{
	//	Message: "Success: SOL Calibration of Port 1 complete",
	//	Command: rc,
	// }

}
