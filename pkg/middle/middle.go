// Package middle coordinates the response to user requests that require the use of the rfswitch and calibration services
package middle

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/timdrysdale/pocket-vna-two-port/pkg/calibration"
	"github.com/timdrysdale/pocket-vna-two-port/pkg/pocket"
	"github.com/timdrysdale/pocket-vna-two-port/pkg/rfswitch"
	"github.com/timdrysdale/pocket-vna-two-port/pkg/stream"
)

type Middle struct {
	Calibration *calibration.Calibration
	Stream      *stream.Stream
	Switch      *rfswitch.Switch
	VNA         *pocket.VNAService
}

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

				case s.Response <- HandleRequest(request, &c, &r, &v):
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
		Calibration: &c,
		Stream:      &s,
		Switch:      &r,
		VNA:         &v,
	}

}

func TimeoutMessage(request interface{}) interface{} {

	return pocket.CustomResult{
		Message: "timeout waiting for request to be handled",
		Command: request,
	}

}

func HandleRequest(request interface{}, c *calibration.Calibration, r *rfswitch.Switch, v *pocket.VNAService) interface{} {

	switch request.(type) {

	case pocket.ReasonableFrequencyRange, pocket.SingleQuery:

		v.Request <- request

		return <-v.Response

	case pocket.RangeQuery:

		// this type is used for different commands

		rq := request.(pocket.RangeQuery)

		switch rq.Command.Command {

		case "rq", "rangequery":

			log.WithFields(log.Fields{
				"request": rq,
			}).Infof("Middle.HandleRequest with ID: %s", rq.ID)

			return RangeQuery(rq, r, v)

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

func CalibratedRangeQuery(crq pocket.CalibratedRangeQuery, c *calibration.Calibration, r *rfswitch.Switch, v *pocket.VNAService) interface{} {

	sc, ok := (c.Scan).(pocket.RangeQuery)

	if !(ok && c.Ready) {
		return pocket.CustomResult{
			Message: "Error. No existing calibration. Please calibrate with rc command",
			Command: crq,
		}
	}

	var err error
	var name string

	switch {
	case crq.What == "short" || crq.What == "s":
		name = "short"
		err = r.SetShort()

	case crq.What == "open" || crq.What == "o":
		name = "open"
		err = r.SetOpen()

	case crq.What == "load" || crq.What == "l":
		name = "load"
		err = r.SetLoad()

	case crq.What == "thru" || crq.What == "t":
		name = "thru"
		err = r.SetThru()

	case crq.What == "dut1" || crq.What == "1":
		name = "dut1"
		err = r.SetDUT1()

	case crq.What == "dut2" || crq.What == "2":
		name = "dut2"
		err = r.SetDUT2()

	case crq.What == "dut3" || crq.What == "3":
		name = "dut3"
		err = r.SetDUT3()

	case crq.What == "dut4" || crq.What == "4":
		name = "dut4"
		err = r.SetDUT4()
	default:
		name = crq.What
		err = fmt.Errorf("unrecognised value of what: %s", name)
	}

	if err != nil {
		return pocket.CustomResult{
			Message: "Error setting RF switch to " + name + ": " + err.Error(),
			Command: crq,
		}
	}

	// modify the scan command to select only
	// the sparams specified by the user's
	// crq command. The cal scans had to do all four sparams,
	// the user might not always want all four
	// the calibration routine does not need four params in the dut
	// to work, according to testing in python of TwelveTerm (TDD Nov 2022)
	sc.Select = crq.Select

	v.Request <- sc

	log.Debugf("Scan request %v", sc)

	response := <-v.Response

	log.Debugf("Scan response %v", response)

	rrq, ok := response.(pocket.RangeQuery)

	log.Debugf("Scan response as range query %v", rrq)

	if !ok {
		return pocket.CustomResult{
			Message: "Error measuring " + name,
			Command: response,
		}
	}

	result := rrq.Result

	if len(result) != sc.Size {
		return pocket.CustomResult{
			Message: "Error measuring " + name,
			Command: response,
		}
	}

	err = c.SetDUTParam(result)

	if err != nil {
		return pocket.CustomResult{
			Message: "Error putting data for " + name + " into cal store as DUT: " + err.Error(),
			Command: result,
		}
	}

	// apply calibration to DUT data
	calibrated, err := c.Apply()

	if err != nil {
		return pocket.CustomResult{
			Message: "Error applying calibration to measured data for " + name + ": " + err.Error(),
			// don't include result - not in correct format and will be nil anyway
		}
	}

	sparams, err := calibration.CalibrationToPocket(calibrated)

	if err != nil {
		return pocket.CustomResult{
			Message: "Error converting calibrated data format for " + name + ": " + err.Error(),
			// don't include result - not in correct format and will be nil anyway
		}
	}

	crq.Result = sparams

	return crq

}

func RangeQuery(rq pocket.RangeQuery, r *rfswitch.Switch, v *pocket.VNAService) interface{} {

	var err error
	var name string

	switch {
	case rq.What == "short" || rq.What == "s":
		name = "short"
		err = r.SetShort()

	case rq.What == "open" || rq.What == "o":
		name = "open"
		err = r.SetOpen()

	case rq.What == "load" || rq.What == "l":
		name = "load"
		err = r.SetLoad()

	case rq.What == "thru" || rq.What == "t":
		name = "thru"
		err = r.SetThru()

	case rq.What == "dut1" || rq.What == "1":
		name = "dut1"
		err = r.SetDUT1()

	case rq.What == "dut2" || rq.What == "2":
		name = "dut2"
		err = r.SetDUT2()

	case rq.What == "dut3" || rq.What == "3":
		name = "dut3"
		err = r.SetDUT3()

	case rq.What == "dut4" || rq.What == "4":
		name = "dut4"
		err = r.SetDUT4()
	}

	// throw no error if what is unrecognised, because it will be blank when rq is used by rangecal and calibratedrangequery
	// ideally we'd use this in the same way for all uses, but using rq externally only became necessary for troubleshooting
	// the two-port rig with 8-port switches, so we do it this way to minimise changes elsewhere for now.
	// but do throw error if the value is what is valid
	if err != nil {
		return pocket.CustomResult{
			Message: "Error setting RF switch to " + name + ": " + err.Error(),
			Command: rq,
		}
	}

	v.Request <- rq

	log.Debugf("Scan request %v", rq)

	response := <-v.Response

	log.Debugf("Scan response %v", response)

	rrq, ok := response.(pocket.RangeQuery)

	log.Debugf("Scan response as range query %v", rrq)

	if !ok {
		return pocket.CustomResult{
			Message: "Error measuring " + name,
			Command: response,
		}
	}

	return rrq

}

func RangeCal(rc pocket.RangeQuery, c *calibration.Calibration, r *rfswitch.Switch, v *pocket.VNAService) interface{} {

	// clear previous cal
	c.Clear()

	// prepare the scanning command used to measure each standard
	scan := rc
	scan.Command.Command = "rq"

	//save it for the cqr to use later
	c.Scan = scan

	// SHORT

	name := "short"

	scan.Select = pocket.SParamSelect{
		S11: true,
		S22: true,
	}

	log.Debugf("Middle.RangeCal: setting rfswitch to %s", name)

	err := r.SetShort()

	if err != nil {
		log.Errorf("Middle.RangeCal error setting %s was %s", name, err.Error())
		return pocket.CustomResult{
			Message: "Error setting RF switch to " + name + ": " + err.Error(),
			Command: rc,
		}
	} else {
		log.Debug("Middle.RangeCal set short ok")
	}

	log.Debug("Middle.RangeCal requesting scan from VNA")

	v.Request <- scan

	log.Debug("Middle.RangeCal awaiting result from VNA")

	response := <-v.Response

	log.Debug("Middle.RangeCal checking result from VNA")

	log.Debugf("response: %s", response)

	rrq, ok := response.(pocket.RangeQuery)

	if !ok {
		log.Errorf("Middle.RangeCal error with scanning %s", name)
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
	scan.Select = pocket.SParamSelect{
		S11: true,
		S22: true,
	}
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
	scan.Select = pocket.SParamSelect{
		S11: true,
		S22: true,
	}
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

	// THRU
	name = "thru"
	scan.Select = pocket.SParamSelect{
		S11: true,
		S12: true,
		S21: true,
		S22: true,
	}
	err = r.SetThru()

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

	err = c.SetThruParam(result)

	if err != nil {
		return pocket.CustomResult{
			Message: "Error putting data for " + name + " into cal store: " + err.Error(),
			Command: result,
		}
	}

	c.Ready = true

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
