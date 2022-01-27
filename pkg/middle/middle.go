// Package middle coordinates the response to user requests that require the use of the rfswitch and calibration services
package middle

import (
	"context"
	"time"

	"github.com/timdrysdale/go-pocketvna/pkg/calibration"
	"github.com/timdrysdale/go-pocketvna/pkg/pocket"
	"github.com/timdrysdale/go-pocketvna/pkg/rfswitch"
	"github.com/timdrysdale/go-pocketvna/pkg/stream"
)

func New(uc, ur, us string, ctx context.Context) Middle {

	// placeholders during testing
	//c := calibration.Calibration{}
	//r := rfswitch.Switch{}
	//s := stream.Stream{}
	//v := pocket.VNAService{}

	c := calibration.New(uc, ctx)
	r := rfswitch.New(ur, ctx)
	s := stream.New(us, ctx)

	v := pocket.New(ctx)

	// avoid compile warnings during testing when services commented out
	//fmt.Printf("cal:    %s\nswitch: %s\nstream: %s\n", uc, ur, us)

	timeout := time.Second

	requesttimeout := 2 * time.Minute

	go func() {
		for {

			select {

			case request := <-s.Request:

				//fmt.Printf("Handling Request - STARTING")

				select {

				case s.Response <- HandleRequest(request, c, r, v):
					// carry on
					//fmt.Printf("Handling Request - DONE")
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

	case pocket.ReasonableFrequencyRange:

		v.Request <- request.(pocket.ReasonableFrequencyRange)

		return <-v.Response

	case pocket.RangeQuery:

		result, err := v.VNA.RangeQuery(request.(pocket.RangeQuery))

		if err != nil {
			return pocket.CustomResult{Message: err.Error()}
		}

		return result

	case pocket.SingleQuery:

		result, err := v.VNA.SingleQuery(request.(pocket.SingleQuery))

		if err != nil {
			return pocket.CustomResult{Message: err.Error()}
		}

		return result

	default:
		return pocket.CustomResult{
			Message: "Unknown request",
			Command: request,
		}
	}
}
